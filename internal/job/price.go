package job

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/cache"
	"github.com/pagu-project/pagu/pkg/log"
)

const (
	_defaultXeggexPriceEndpoint = "https://api.xeggex.com/api/v2/market/getbysymbol/Pactus%2Fusdt"
	_defaultTradeOgreEndpoint   = "https://tradeogre.com/api/v1/ticker/PAC-USDT"
	_defaultAzbitPriceEndpoint  = "https://data.azbit.com/api/tickers?currencyPairCode=PAC_USDT"
)

type price struct {
	ctx    context.Context
	cache  cache.Cache[string, entity.Price]
	ticker *time.Ticker
	cancel context.CancelFunc
}

func NewPrice(
	cch cache.Cache[string, entity.Price],
) Job {
	ctx, cancel := context.WithCancel(context.Background())

	return &price{
		cache:  cch,
		ticker: time.NewTicker(128 * time.Second),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *price) Start() {
	p.start()
	go p.runTicker()
}

func (p *price) start() {
	var (
		wg        sync.WaitGroup
		price     entity.Price
		xeggex    entity.XeggexPriceResponse
		tradeOgre entity.TradeOgrePriceResponse
		azbit     []entity.AzbitPriceResponse
	)

	ctx := context.Background()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.getPrice(ctx, _defaultXeggexPriceEndpoint, &xeggex); err != nil {
			log.Error(err.Error())

			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.getPrice(ctx, _defaultTradeOgreEndpoint, &tradeOgre); err != nil {
			log.Error(err.Error())

			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.getPrice(ctx, _defaultAzbitPriceEndpoint, &azbit); err != nil {
			log.Error(err.Error())

			return
		}
	}()

	wg.Wait()

	price.XeggexPacToUSDT = xeggex
	price.TradeOgrePacToUSDT = tradeOgre
	if len(azbit) > 0 {
		price.AzbitPacToUSDT = azbit[0]
	}

	ok := p.cache.Exists(config.PriceCacheKey)
	if ok {
		p.cache.Update(config.PriceCacheKey, price, 0)
	} else {
		p.cache.Add(config.PriceCacheKey, price, 0)
	}
}

func (p *price) runTicker() {
	for {
		select {
		case <-p.ctx.Done():
			return

		case <-p.ticker.C:
			p.start()
		}
	}
}

func (*price) getPrice(ctx context.Context, endpoint string, priceResponse any) error {
	cli := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response code is %v", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)

	return dec.Decode(priceResponse)
}

func (p *price) Stop() {
	p.ticker.Stop()
}
