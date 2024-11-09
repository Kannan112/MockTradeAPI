package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/request"
	"github.com/kannan112/mock-trading-platform-api/pkg/api/handler/response"
	"github.com/kannan112/mock-trading-platform-api/pkg/repository/interfaces"
	"github.com/kannan112/mock-trading-platform-api/pkg/service/token"
	service "github.com/kannan112/mock-trading-platform-api/pkg/usecase/interfaces"
	"github.com/kannan112/mock-trading-platform-api/pkg/utils"
)

type userUserCase struct {
	userRepo     interfaces.UserRepository
	orderRepo    interfaces.OrderRepository
	tokenService token.TokenService
}

func NewUserUseCase(userRepo interfaces.UserRepository, tokenService token.TokenService, orderRepo interfaces.OrderRepository) service.UserUseCase {
	return &userUserCase{
		userRepo:     userRepo,
		orderRepo:    orderRepo,
		tokenService: tokenService,
	}
}

const (
	AccessTokenDuration = time.Minute * 20
)

func (c *userUserCase) CeateNewUser(ctx context.Context, body request.RegisterUserRequest) error {

	exists, err := c.userRepo.FindUserByEmail(ctx, body.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	password, err := utils.GenerateHashFromPassword(body.Password)
	if err != nil {
		return err
	}

	body.Password = password

	_, err = c.userRepo.SaveUser(ctx, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *userUserCase) UserLogin(ctx context.Context, body request.LoginRequest) (response.Token, error) {

	exists, err := c.userRepo.FindUserByEmail(ctx, body.Email)
	if err != nil {
		return response.Token{}, err
	}

	if !exists {
		return response.Token{}, errors.New("user not exists")
	}

	hashPassword, err := c.userRepo.ExtractPassword(ctx, body.Email)
	if err != nil {
		return response.Token{}, err
	}

	verify := utils.VerifyHashAndPassword(hashPassword, body.Password)
	if !verify {
		return response.Token{}, errors.New("worng password")
	}

	uid, err := c.userRepo.GetUserId(ctx, body.Email)
	if err != nil {
		return response.Token{}, err
	}

	token, err := token.GenerateAccessToken(uid)

	if err != nil {
		return response.Token{}, err
	}

	return response.Token{
		AccessToken: token,
	}, nil

}

func (c *userUserCase) FetchMarketData(symbol string) (response.MarketData, error) {
	formattedSymbol := strings.ToUpper(strings.ReplaceAll(symbol, " ", ""))

	// Add USDT if not present (assuming USDT is the default quote currency)
	if !strings.HasSuffix(formattedSymbol, "USDT") {
		formattedSymbol = formattedSymbol + "USDT"
	}

	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/bookTicker?symbol=%s", formattedSymbol)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Log the URL we're calling (for debugging)
	log.Printf("Calling Binance API with URL: %s", url)

	resp, err := client.Get(url)
	if err != nil {
		return response.MarketData{}, fmt.Errorf("failed to fetch market data: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.MarketData{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var binanceError struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}
		if err := json.Unmarshal(bodyBytes, &binanceError); err != nil {
			return response.MarketData{}, fmt.Errorf("API error: %s", string(bodyBytes))
		}
		return response.MarketData{}, fmt.Errorf("Binance API error: %s (code: %d)",
			binanceError.Msg, binanceError.Code)
	}

	var data response.MarketData
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return response.MarketData{}, fmt.Errorf("failed to decode market data: %w, raw data: %s",
			err, string(bodyBytes))
	}

	if data.BidPrice == 0 || data.AskPrice == 0 {
		return response.MarketData{}, fmt.Errorf("received invalid price data for symbol %s: bid=%v, ask=%v",
			formattedSymbol, data.BidPrice, data.AskPrice)
	}

	return data, nil
}

// Helper function to validate if a symbol is supported by Binance
func isValidBinanceSymbol(symbol string) bool {
	validSymbols := map[string]bool{
		"BTCUSDT":  true,
		"ETHUSDT":  true,
		"BNBUSDT":  true,
		"ADAUSDT":  true,
		"DOGEUSDT": true,
		"XRPUSDT":  true,
		// Add more symbols as needed
	}

	return validSymbols[strings.ToUpper(symbol)]
}
func (c *userUserCase) GetMarketPrice(marketData response.MarketData, orderType string) (float64, error) {
	// Convert orderType to lowercase for case-insensitive comparison
	orderType = strings.ToLower(orderType)

	switch orderType {
	case "buy":
		if marketData.AskPrice <= 0 {
			return 0, fmt.Errorf("invalid ask price: %v", marketData.AskPrice)
		}
		return marketData.AskPrice, nil
	case "sell":
		if marketData.BidPrice <= 0 {
			return 0, fmt.Errorf("invalid bid price: %v", marketData.BidPrice)
		}
		return marketData.BidPrice, nil
	default:
		return 0, fmt.Errorf("invalid order type: %s", orderType)
	}
}

func (c *userUserCase) CreateOrder(ctx context.Context, uid int, orderData response.OrderResponse) (oid int, err error) {

	oid, err = c.orderRepo.PlaceOrder(ctx, uid, orderData)

	return oid, err
}

func (c *userUserCase) ListOrders(uid int) ([]utils.OrderResponse, error) {
	data, err := c.orderRepo.GetAllOrders(uid)
	return data, err
}

func (c *userUserCase) GetOrderByID(ctx context.Context, uid, oid uint) (utils.Order, error) {
	data, err := c.orderRepo.GetOrderByID(oid, uid)
	return data, err
}

func (c *userUserCase) DeleteOrderById(ctx context.Context, uid, oid uint) error {
	err := c.orderRepo.DeleteOrderById(oid, uid)
	return err
}
