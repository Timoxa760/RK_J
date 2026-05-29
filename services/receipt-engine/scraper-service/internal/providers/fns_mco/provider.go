package fns_mco

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jfk9w-go/based"
	"github.com/jfk9w-go/lkdr-api"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type Provider struct {
	ruCaptchaKey string
	tokenDir     string
	httpClient   *http.Client

	mu    sync.Mutex
	auths map[string]*authFlow
}

type authFlow struct {
	codeCh chan string
	done   chan error
}

func NewProvider(ruCaptchaKey, tokenDir string) *Provider {
	if tokenDir == "" {
		tokenDir = "mco_tokens"
	}
	os.MkdirAll(tokenDir, 0755)
	return &Provider{
		ruCaptchaKey: ruCaptchaKey,
		tokenDir:     tokenDir,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		auths:        make(map[string]*authFlow),
	}
}

func (p *Provider) StartAuth(phone string) error {
	ts := &fileTokenStorage{dir: p.tokenDir}

	tokens, err := ts.LoadTokens(context.Background(), phone)
	if err != nil {
		return fmt.Errorf("check tokens: %w", err)
	}
	if tokens != nil {
		return nil
	}

	p.mu.Lock()
	if _, ok := p.auths[phone]; ok {
		p.mu.Unlock()
		return fmt.Errorf("auth already in progress for %s", phone)
	}

	af := &authFlow{
		codeCh: make(chan string, 1),
		done:   make(chan error, 1),
	}
	p.auths[phone] = af
	p.mu.Unlock()

	client, err := lkdr.NewClient(lkdr.ClientParams{
		Phone:        phone,
		Clock:        based.StandardClock,
		DeviceID:     fmt.Sprintf("mco-web-%d", time.Now().UnixNano()%1000000),
		UserAgent:    defaultUserAgent,
		TokenStorage: ts,
	})
	if err != nil {
		p.mu.Lock()
		delete(p.auths, phone)
		p.mu.Unlock()
		return fmt.Errorf("create lkdr client: %w", err)
	}

	auth := &ruCaptchaAuthorizer{
		ruCaptchaKey: p.ruCaptchaKey,
		httpClient:   p.httpClient,
		codeCh:       af.codeCh,
	}
	authCtx := lkdr.WithAuthorizer(context.Background(), auth)

	go func() {
		_, err := client.Receipt(authCtx, &lkdr.ReceiptIn{Limit: 1, Offset: 0})
		af.done <- err

		p.mu.Lock()
		delete(p.auths, phone)
		p.mu.Unlock()
	}()

	return nil
}

func (p *Provider) VerifyAuth(phone, code string) error {
	p.mu.Lock()
	af, ok := p.auths[phone]
	p.mu.Unlock()

	if !ok {
		ts := &fileTokenStorage{dir: p.tokenDir}
		tokens, err := ts.LoadTokens(context.Background(), phone)
		if err != nil {
			return fmt.Errorf("check tokens: %w", err)
		}
		if tokens != nil {
			return nil
		}
		return fmt.Errorf("no pending auth for %s; call /auth first", phone)
	}

	select {
	case af.codeCh <- code:
	default:
		return fmt.Errorf("code channel full, auth may be stuck")
	}

	err := <-af.done
	return err
}

func (p *Provider) SyncReceipts(ctx context.Context, phone string) ([]scrap.RawReceipt, error) {
	ts := &fileTokenStorage{dir: p.tokenDir}

	client, err := lkdr.NewClient(lkdr.ClientParams{
		Phone:        phone,
		Clock:        based.StandardClock,
		DeviceID:     fmt.Sprintf("mco-web-%d", time.Now().UnixNano()%1000000),
		UserAgent:    defaultUserAgent,
		TokenStorage: ts,
	})
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	var allReceipts []scrap.RawReceipt
	offset := 0
	limit := 100

	for {
		out, err := client.Receipt(ctx, &lkdr.ReceiptIn{
			Limit:   limit,
			Offset:  offset,
			OrderBy: "receiveDate",
		})
		if err != nil {
			return nil, fmt.Errorf("fetch receipts at offset %d: %w", offset, err)
		}

		for _, r := range out.Receipts {
			if r.Key == "" {
				continue
			}

			fd, err := client.FiscalData(ctx, &lkdr.FiscalDataIn{Key: r.Key})
			if err != nil {
				if lkdr.IsDataNotFound(err) {
					continue
				}
				return nil, fmt.Errorf("fetch fiscal data for %s: %w", r.Key, err)
			}

			raw := mapToRawReceipt(fd, phone)
			allReceipts = append(allReceipts, *raw)
		}

		if !out.HasMore {
			break
		}
		offset += limit
	}

	return allReceipts, nil
}

type fileTokenStorage struct {
	dir string
	mu  sync.Mutex
}

func (s *fileTokenStorage) LoadTokens(_ context.Context, phone string) (*lkdr.Tokens, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(filepath.Join(s.dir, phone+".json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var tokens lkdr.Tokens
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, err
	}
	return &tokens, nil
}

func (s *fileTokenStorage) UpdateTokens(_ context.Context, phone string, tokens *lkdr.Tokens) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.dir, phone+".json")
	data, err := json.Marshal(tokens)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

type ruCaptchaAuthorizer struct {
	ruCaptchaKey string
	httpClient   *http.Client
	codeCh       <-chan string
}

func (a *ruCaptchaAuthorizer) GetCaptchaToken(ctx context.Context, userAgent, siteKey, pageURL string) (string, error) {
	if a.ruCaptchaKey == "" {
		return "", fmt.Errorf("RUCAPTCHA_KEY not set")
	}

	inURL := fmt.Sprintf("http://rucaptcha.com/in.php?key=%s&method=userrecaptcha&googlekey=%s&pageurl=%s",
		a.ruCaptchaKey, siteKey, pageURL)

	req, err := http.NewRequestWithContext(ctx, "GET", inURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("rucaptcha send: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result := strings.TrimSpace(string(body))

	if !strings.HasPrefix(result, "OK|") {
		return "", fmt.Errorf("rucaptcha send failed: %s", result)
	}

	captchaID := strings.TrimPrefix(result, "OK|")
	pollURL := fmt.Sprintf("http://rucaptcha.com/res.php?key=%s&action=get&id=%s", a.ruCaptchaKey, captchaID)

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(3 * time.Second):
		}

		req, err := http.NewRequestWithContext(ctx, "GET", pollURL, nil)
		if err != nil {
			return "", err
		}

		resp, err := a.httpClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("rucaptcha poll: %w", err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		result = strings.TrimSpace(string(body))

		if result == "CAPCHA_NOT_READY" {
			continue
		}

		if !strings.HasPrefix(result, "OK|") {
			return "", fmt.Errorf("rucaptcha poll failed: %s", result)
		}

		return strings.TrimPrefix(result, "OK|"), nil
	}
}

func (a *ruCaptchaAuthorizer) GetConfirmationCode(ctx context.Context, _ string) (string, error) {
	select {
	case code := <-a.codeCh:
		return code, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func mapToRawReceipt(fd *lkdr.FiscalDataOut, userID string) *scrap.RawReceipt {
	store := ""
	if fd.RetailPlace != nil {
		store = *fd.RetailPlace
	}

	rp := &scrap.RawReceipt{
		ID:       fmt.Sprintf("fns-%s-%d", fd.FiscalDriveNumber, fd.FiscalDocumentNumber),
		UserID:   userID,
		Provider: "fns",
		Store:    store,
		Date:     fd.DateTime.Time(),
		Total:    fd.TotalSum,
	}

	for _, item := range fd.Items {
		rp.Items = append(rp.Items, scrap.RawItem{
			Name:     strings.TrimSpace(item.Name),
			Price:    item.Price,
			Quantity: int(item.Quantity),
		})
	}

	return rp
}

func mapToRawReceiptFromFiscalData(fd *lkdr.FiscalDataOut, userID string, fn, fdStr string) *scrap.RawReceipt {
	rp := mapToRawReceipt(fd, userID)
	rp.ID = fmt.Sprintf("fns-%s-%s", fn, fdStr)
	return rp
}
