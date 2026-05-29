package email

import (
	"context"
	"fmt"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-sasl"
)

type IMAPConfig struct {
	Address  string
	Username string
}

type IMAPClient struct {
	oauthMgr *OAuthManager
	configs  map[Provider]IMAPConfig
	demoMode bool
}

func NewIMAPClient(oauthMgr *OAuthManager) *IMAPClient {
	demoMode := oauthMgr.demoMode

	return &IMAPClient{
		oauthMgr: oauthMgr,
		demoMode: demoMode,
		configs: map[Provider]IMAPConfig{
			ProviderYandex: {Address: "imap.yandex.ru:993", Username: ""},
			ProviderMail:   {Address: "imap.mail.ru:993", Username: ""},
		},
	}
}

func (c *IMAPClient) FetchReceipts(ctx context.Context, provider Provider, userID string) ([]string, error) {
	if c.demoMode {
		return c.fakeFetch(provider)
	}

	cfg, ok := c.configs[provider]
	if !ok {
		return nil, fmt.Errorf("imap: unknown provider %s", provider)
	}

	tok, err := c.oauthMgr.GetToken(ctx, provider, userID)
	if err != nil {
		return nil, fmt.Errorf("imap: get token: %w", err)
	}

	client, err := imapclient.DialTLS(cfg.Address, nil)
	if err != nil {
		return nil, fmt.Errorf("imap: dial %s: %w", cfg.Address, err)
	}
	defer client.Close()

	saslClient := sasl.NewOAuthBearerClient(&sasl.OAuthBearerOptions{
		Username: cfg.Username,
		Token:    tok.AccessToken,
	})
	if err := client.Authenticate(saslClient); err != nil {
		return nil, fmt.Errorf("imap: auth: %w", err)
	}

	selectCmd := client.Select("INBOX", nil)
	_, err = selectCmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("imap: select INBOX: %w", err)
	}

	searchCmd := client.Search(&imap.SearchCriteria{
		Since: time.Now().AddDate(0, 0, -30),
		Or: [][2]imap.SearchCriteria{
			{
				{Header: []imap.SearchCriteriaHeaderField{{Key: "SUBJECT", Value: "чек"}}},
				{Header: []imap.SearchCriteriaHeaderField{{Key: "SUBJECT", Value: "ticket"}}},
			},
			{
				{Header: []imap.SearchCriteriaHeaderField{{Key: "SUBJECT", Value: "чека"}}},
				{Header: []imap.SearchCriteriaHeaderField{{Key: "SUBJECT", Value: "receipt"}}},
			},
		},
	}, nil)

	searchData, err := searchCmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("imap: search: %w", err)
	}

	seqNums := searchData.AllSeqNums()
	if len(seqNums) == 0 {
		return nil, nil
	}

	fetchCmd := client.Fetch(imap.SeqSetNum(seqNums...), &imap.FetchOptions{
		BodySection: []*imap.FetchItemBodySection{
			{},
		},
	})

	messages, err := fetchCmd.Collect()
	if err != nil {
		return nil, fmt.Errorf("imap: fetch: %w", err)
	}

	var results []string
	for _, msg := range messages {
		for _, section := range msg.BodySection {
			results = append(results, string(section.Bytes))
		}
	}

	return results, nil
}

func (c *IMAPClient) fakeFetch(provider Provider) ([]string, error) {
	return []string{
		fmt.Sprintf(`<html><body><table>
			<tr><td>Магазин</td><td>Пятёрочка</td></tr>
			<tr><td>Дата</td><td>%s</td></tr>
			<tr><td>Сумма</td><td>1 032.50 ₽</td></tr>
			<tr><td>Молоко 3.2%%</td><td>78.00 ₽</td><td>1</td></tr>
			<tr><td>Хлеб белый</td><td>45.00 ₽</td><td>1</td></tr>
			<tr><td>Сыр Российский</td><td>189.00 ₽</td><td>1</td></tr>
		</table></body></html>`, time.Now().Format("2006-01-02")),
		fmt.Sprintf(`<html><body><table>
			<tr><td>Магазин</td><td>ВкусВилл</td></tr>
			<tr><td>Дата</td><td>%s</td></tr>
			<tr><td>Сумма</td><td>654.00 ₽</td></tr>
			<tr><td>Творог 5%%</td><td>120.00 ₽</td><td>1</td></tr>
			<tr><td>Кефир</td><td>85.00 ₽</td><td>1</td></tr>
		</table></body></html>`, time.Now().AddDate(0, 0, -1).Format("2006-01-02")),
	}, nil
}
