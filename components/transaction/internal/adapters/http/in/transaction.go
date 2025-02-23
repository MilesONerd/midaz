package in

import (
	"context"
	"encoding/json"
	"github.com/LerianStudio/midaz/components/transaction/internal/adapters/postgres/operation"
	"github.com/LerianStudio/midaz/components/transaction/internal/adapters/postgres/transaction"
	"github.com/LerianStudio/midaz/components/transaction/internal/services/command"
	"github.com/LerianStudio/midaz/components/transaction/internal/services/query"
	"github.com/LerianStudio/midaz/pkg"
	"github.com/LerianStudio/midaz/pkg/constant"
	goldTransaction "github.com/LerianStudio/midaz/pkg/gold/transaction"
	goldModel "github.com/LerianStudio/midaz/pkg/gold/transaction/model"
	"github.com/LerianStudio/midaz/pkg/mgrpc/account"
	"github.com/LerianStudio/midaz/pkg/mlog"
	"github.com/LerianStudio/midaz/pkg/mmodel"
	"github.com/LerianStudio/midaz/pkg/mopentelemetry"
	"github.com/LerianStudio/midaz/pkg/mpostgres"
	"github.com/LerianStudio/midaz/pkg/net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/otel/trace"
	"os"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TransactionHandler struct that handle transaction
type TransactionHandler struct {
	Command *command.UseCase
	Query   *query.UseCase
}

// CreateTransactionJSON method that create transaction using JSON
//
//	@Summary		Create a Transaction using JSON
//	@Description	Create a Transaction with the input payload
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Authorization Bearer Token"
//	@Param			Midaz-Id		header		string								false	"Request ID"
//	@Param			organization_id	path		string								true	"Organization ID"
//	@Param			ledger_id		path		string								true	"Ledger ID"
//	@Param			transaction		body		transaction.CreateTransactionInput	true	"Transaction Input"
//	@Success		200				{object}	transaction.Transaction
//	@Router			/v1/organizations/{organization_id}/ledgers/{ledger_id}/transactions/json [post]
func (handler *TransactionHandler) CreateTransactionJSON(p any, c *fiber.Ctx) error {
	ctx := c.UserContext()

	logger := pkg.NewLoggerFromContext(ctx)
	tracer := pkg.NewTracerFromContext(ctx)

	ctx, span := tracer.Start(ctx, "handler.create_transaction")
	defer span.End()

	c.SetUserContext(ctx)

	input := p.(*transaction.CreateTransactionInput)
	parserDSL := input.FromDSl()
	logger.Infof("Request to create an transaction with details: %#v", parserDSL)

	response := handler.createTransaction(c, logger, *parserDSL)

	return response
}

// CreateTransactionDSL method that create transaction using DSL
//
//	@Summary		Create a Transaction using DSL
//	@Description	Create a Transaction with the input DSL file
//	@Tags			Transactions
//	@Accept			mpfd
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authorization Bearer Token"
//	@Param			Midaz-Id		header		string	false	"Request ID"
//	@Param			organization_id	path		string	true	"Organization ID"
//	@Param			ledger_id		path		string	true	"Ledger ID"
//	@Param			transaction		formData	file	true	"Transaction DSL file"
//	@Success		200				{object}	transaction.Transaction
//	@Router			/v1/organizations/{organization_id}/ledgers/{ledger_id}/transactions/dsl [post]
func (handler *TransactionHandler) CreateTransactionDSL(c *fiber.Ctx) error {
	ctx := c.UserContext()

	logger := pkg.NewLoggerFromContext(ctx)
	tracer := pkg.NewTracerFromContext(ctx)

	ctx, span := tracer.Start(ctx, "handler.create_transaction_dsl")
	defer span.End()

	c.SetUserContext(ctx)

	_, err := http.ValidateParameters(c.Queries())
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to validate query parameters", err)

		logger.Error("Failed to validate query parameters: ", err.Error())

		return http.WithError(c, err)
	}

	dsl, err := http.GetFileFromHeader(c)
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to get file from Header", err)

		logger.Error("Failed to get file from Header: ", err.Error())

		return http.WithError(c, err)
	}

	errListener := goldTransaction.Validate(dsl)
	if errListener != nil && len(errListener.Errors) > 0 {
		err := pkg.ValidateBusinessError(constant.ErrInvalidDSLFileFormat, reflect.TypeOf(transaction.Transaction{}).Name(), errListener.Errors)

		mopentelemetry.HandleSpanError(&span, "Failed to validate script in DSL", err)

		return http.WithError(c, err)
	}

	parsed := goldTransaction.Parse(dsl)

	parserDSL, ok := parsed.(goldModel.Transaction)
	if !ok {
		err := pkg.ValidateBusinessError(constant.ErrInvalidDSLFileFormat, reflect.TypeOf(transaction.Transaction{}).Name())

		mopentelemetry.HandleSpanError(&span, "Failed to parse script in DSL", err)

		return http.WithError(c, err)
	}

	response := handler.createTransaction(c, logger, parserDSL)

	return response
}

// CreateTransactionTemplate method that create transaction template
//
// TODO: Implement this method and the swagger documentation related to it
func (handler *TransactionHandler) CreateTransactionTemplate(p any, c *fiber.Ctx) error {
	ctx := c.UserContext()

	logger := pkg.NewLoggerFromContext(ctx)
	tracer := pkg.NewTracerFromContext(ctx)

	tracer.Start(ctx, "handler.create_transaction_template")

	payload := p.(*transaction.InputDSL)
	logger.Infof("Request to create an transaction with details: %#v", payload)

	return http.Created(c, payload)
}

// CommitTransaction method that commit transaction created before
//
// TODO: Implement this method and the swagger documentation related to it
func (handler *TransactionHandler) CommitTransaction(c *fiber.Ctx) error {
	ctx := c.UserContext()

	logger := pkg.NewLoggerFromContext(ctx)
	tracer := pkg.NewTracerFromContext(ctx)

	_, span := tracer.Start(ctx, "handler.commit_transaction")
	defer span.End()

	return http.Created(c, logger)
}

// RevertTransaction method that revert transaction created before
//
//	@Summary		Revert a Transaction
//	@Description	Revert a Transaction with Transaction ID only
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Authorization Bearer Token"
//	@Param			Midaz-Id		header		string								false	"Request ID"
//	@Param			organization_id	path		string								true	"Organization ID"
//	@Param			ledger_id		path		string								true	"Ledger ID"
//	@Param			transaction_id	path		string								true	"Transaction ID"
//	@Success		200				{object}	transaction.Transaction
//	@Router			/v1/organizations/{organization_id}/ledgers/{ledger_id}/transactions/{transaction_id}/revert [post]
func (handler *TransactionHandler) RevertTransaction(c *fiber.Ctx) error {
	ctx := c.UserContext()

	logger := pkg.NewLoggerFromContext(ctx)
	tracer := pkg.NewTracerFromContext(ctx)

	_, span := tracer.Start(ctx, "handler.revert_transaction")
	defer span.End()

	organizationID := c.Locals("organization_id").(uuid.UUID)
	ledgerID := c.Locals("ledger_id").(uuid.UUID)
	transactionID := c.Locals("transaction_id").(uuid.UUID)

	parent, err := handler.Query.GetParentByTransactionID(ctx, organizationID, ledgerID, transactionID)
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to retrieve Parent Transaction on query", err)

		logger.Errorf("Failed to retrieve Parent Transaction with ID: %s, Error: %s", transactionID.String(), err.Error())

		return http.WithError(c, err)
	}

	if parent != nil {
		err = pkg.ValidateBusinessError(constant.ErrTransactionIDHasAlreadyParentTransaction, "RevertTransaction")

		mopentelemetry.HandleSpanError(&span, "Transaction Has Already Parent Transaction", err)

		logger.Errorf("Transaction Has Already Parent Transaction with ID: %s, Error: %s", transactionID.String(), err)

		return http.WithError(c, err)
	}

	tran, err := handler.Query.GetTransactionByID(ctx, organizationID, ledgerID, transactionID)
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to retrieve transaction on query", err)

		logger.Errorf("Failed to retrieve Transaction with ID: %s, Error: %s", transactionID.String(), err.Error())

		return http.WithError(c, err)
	}

	if tran.ParentTransactionID != nil {
		err = pkg.ValidateBusinessError(constant.ErrTransactionIDIsAlreadyARevert, "RevertTransaction")

		mopentelemetry.HandleSpanError(&span, "Transaction Has Already Parent Transaction", err)

		logger.Errorf("Transaction Has Already Parent Transaction with ID: %s, Error: %s", transactionID.String(), err)

		return http.WithError(c, err)
	}

	transactionReverted := tran.TransactionRevert()
	if transactionReverted.IsEmpty() {
		err = pkg.ValidateBusinessError(constant.ErrTransactionCantRevert, "RevertTransaction")

		mopentelemetry.HandleSpanError(&span, "Transaction can't be reverted", err)

		logger.Errorf("Parent Transaction can't be reverted with ID: %s, Error: %s", transactionID.String(), err)

		return http.WithError(c, err)
	}

	response := handler.createTransaction(c, logger, transactionReverted)

	return response
}

// UpdateTransaction method that patch transaction created before
//
//	@Summary		Update a Transaction
//	@Description	Update a Transaction with the input payload
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Authorization Bearer Token"
//	@Param			Midaz-Id		header		string								false	"Request ID"
//	@Param			organization_id	path		string								true	"Organization ID"
//	@Param			ledger_id		path		string								true	"Ledger ID"
//	@Param			transaction_id	path		string								true	"Transaction ID"
//	@Param			transaction		body		transaction.UpdateTransactionInput	true	"Transaction Input"
//	@Success		200				{object}	transaction.Transaction
//	@Router			/v1/organizations/{organization_id}/ledgers/{ledger_id}/transactions/{transaction_id} [patch]
func (handler *TransactionHandler) UpdateTransaction(p any, c *fiber.Ctx) error {
	ctx := c.UserContext()

	logger := pkg.NewLoggerFromContext(ctx)
	tracer := pkg.NewTracerFromContext(ctx)

	ctx, span := tracer.Start(ctx, "handler.update_transaction")
	defer span.End()

	organizationID := c.Locals("organization_id").(uuid.UUID)
	ledgerID := c.Locals("ledger_id").(uuid.UUID)
	transactionID := c.Locals("transaction_id").(uuid.UUID)

	logger.Infof("Initiating update of Transaction with Organization ID: %s, Ledger ID: %s and ID: %s", organizationID.String(), ledgerID.String(), transactionID.String())

	payload := p.(*transaction.UpdateTransactionInput)
	logger.Infof("Request to update an Transaction with details: %#v", payload)

	_, err := handler.Command.UpdateTransaction(ctx, organizationID, ledgerID, transactionID, payload)
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to update transaction on command", err)

		logger.Errorf("Failed to update Transaction with ID: %s, Error: %s", transactionID.String(), err.Error())

		return http.WithError(c, err)
	}

	trans, err := handler.Query.GetTransactionByID(ctx, organizationID, ledgerID, transactionID)
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to retrieve transaction on query", err)

		logger.Errorf("Failed to retrieve Transaction with ID: %s, Error: %s", transactionID.String(), err.Error())

		return http.WithError(c, err)
	}

	logger.Infof("Successfully updated Transaction with Organization ID: %s, Ledger ID: %s and ID: %s", organizationID.String(), ledgerID.String(), transactionID.String())

	return http.OK(c, trans)
}

// GetTransaction method that get transaction created before
//
//	@Summary		Get a Transaction by ID
//	@Description	Get a Transaction with the input ID
//	@Tags			Transactions
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authorization Bearer Token"
//	@Param			Midaz-Id		header		string	false	"Request ID"
//	@Param			organization_id	path		string	true	"Organization ID"
//	@Param			ledger_id		path		string	true	"Ledger ID"
//	@Param			transaction_id	path		string	true	"Transaction ID"
//	@Success		200				{object}	transaction.Transaction
//	@Router			/v1/organizations/{organization_id}/ledgers/{ledger_id}/transactions/{transaction_id} [get]
func (handler *TransactionHandler) GetTransaction(c *fiber.Ctx) error {
	ctx := c.UserContext()

	logger := pkg.NewLoggerFromContext(ctx)
	tracer := pkg.NewTracerFromContext(ctx)

	ctx, span := tracer.Start(ctx, "handler.get_transaction")
	defer span.End()

	organizationID := c.Locals("organization_id").(uuid.UUID)
	ledgerID := c.Locals("ledger_id").(uuid.UUID)
	transactionID := c.Locals("transaction_id").(uuid.UUID)

	tran, err := handler.Query.GetTransactionByID(ctx, organizationID, ledgerID, transactionID)
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to retrieve transaction on query", err)

		logger.Errorf("Failed to retrieve Transaction with ID: %s, Error: %s", transactionID.String(), err.Error())

		return http.WithError(c, err)
	}

	ctxGetTransaction, spanGetTransaction := tracer.Start(ctx, "handler.get_transaction.get_operations")

	headerParams, err := http.ValidateParameters(c.Queries())
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to validate query parameters", err)

		logger.Errorf("Failed to validate query parameters, Error: %s", err.Error())

		return http.WithError(c, err)
	}

	headerParams.Metadata = &bson.M{}

	tran, err = handler.Query.GetOperationsByTransaction(ctxGetTransaction, organizationID, ledgerID, tran, *headerParams)
	if err != nil {
		mopentelemetry.HandleSpanError(&spanGetTransaction, "Failed to retrieve Operations", err)

		logger.Errorf("Failed to retrieve Operations with ID: %s, Error: %s", tran.ID, err.Error())

		return http.WithError(c, err)
	}

	spanGetTransaction.End()

	logger.Infof("Successfully retrieved Transaction with ID: %s", transactionID.String())

	return http.OK(c, tran)
}

// GetAllTransactions method that get all transactions created before
//
//	@Summary		Get all Transactions
//	@Description	Get all Transactions with the input metadata or without metadata
//	@Tags			Transactions
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authorization Bearer Token"
//	@Param			Midaz-Id		header		string	false	"Request ID"
//	@Param			organization_id	path		string	true	"Organization ID"
//	@Param			ledger_id		path		string	true	"Ledger ID"
//	@Param			limit			query		int		false	"Limit"			default(10)
//	@Param			start_date		query		string	false	"Start Date"	example "2021-01-01"
//	@Param			end_date		query		string	false	"End Date"		example "2021-01-01"
//	@Param			sort_order		query		string	false	"Sort Order"	Enums(asc,desc)
//	@Param			cursor			query		string	false	"Cursor"
//	@Success		200				{object}	mpostgres.Pagination{items=[]transaction.Transaction,next_cursor=string,prev_cursor=string,limit=int,page=nil}
//	@Router			/v1/organizations/{organization_id}/ledgers/{ledger_id}/transactions [get]
func (handler *TransactionHandler) GetAllTransactions(c *fiber.Ctx) error {
	ctx := c.UserContext()

	logger := pkg.NewLoggerFromContext(ctx)
	tracer := pkg.NewTracerFromContext(ctx)

	ctx, span := tracer.Start(ctx, "handler.get_all_transactions")
	defer span.End()

	organizationID := c.Locals("organization_id").(uuid.UUID)
	ledgerID := c.Locals("ledger_id").(uuid.UUID)

	headerParams, err := http.ValidateParameters(c.Queries())
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to validate query parameters", err)

		logger.Errorf("Failed to validate query parameters, Error: %s", err.Error())

		return http.WithError(c, err)
	}

	pagination := mpostgres.Pagination{
		Limit:      headerParams.Limit,
		NextCursor: headerParams.Cursor,
		SortOrder:  headerParams.SortOrder,
		StartDate:  headerParams.StartDate,
		EndDate:    headerParams.EndDate,
	}

	if headerParams.Metadata != nil {
		logger.Infof("Initiating retrieval of all Transactions by metadata")

		err := mopentelemetry.SetSpanAttributesFromStruct(&span, "headerParams", headerParams)
		if err != nil {
			mopentelemetry.HandleSpanError(&span, "Failed to convert metadata headerParams to JSON string", err)

			return http.WithError(c, err)
		}

		trans, err := handler.Query.GetAllMetadataTransactions(ctx, organizationID, ledgerID, *headerParams)
		if err != nil {
			mopentelemetry.HandleSpanError(&span, "Failed to retrieve all Transactions by metadata", err)

			logger.Errorf("Failed to retrieve all Transactions, Error: %s", err.Error())

			return http.WithError(c, err)
		}

		logger.Infof("Successfully retrieved all Transactions by metadata")

		pagination.SetItems(trans)

		return http.OK(c, pagination)
	}

	logger.Infof("Initiating retrieval of all Transactions ")

	headerParams.Metadata = &bson.M{}

	err = mopentelemetry.SetSpanAttributesFromStruct(&span, "headerParams", headerParams)
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to convert headerParams to JSON string", err)

		return http.WithError(c, err)
	}

	trans, cur, err := handler.Query.GetAllTransactions(ctx, organizationID, ledgerID, *headerParams)
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to retrieve all Transactions", err)

		logger.Errorf("Failed to retrieve all Transactions, Error: %s", err.Error())

		return http.WithError(c, err)
	}

	logger.Infof("Successfully retrieved all Transactions")

	pagination.SetItems(trans)
	pagination.SetCursor(cur.Next, cur.Prev)

	return http.OK(c, pagination)
}

// createTransaction func that received struct from DSL parsed and create Transaction
func (handler *TransactionHandler) createTransaction(c *fiber.Ctx, logger mlog.Logger, parserDSL goldModel.Transaction) error {
	ctx := c.UserContext()
	tracer := pkg.NewTracerFromContext(ctx)

	organizationID := c.Locals("organization_id").(uuid.UUID)
	ledgerID := c.Locals("ledger_id").(uuid.UUID)
	transactionID, _ := c.Locals("transaction_id").(uuid.UUID)

	_, spanIdempotency := tracer.Start(ctx, "handler.create_transaction_idempotency")

	ts, _ := pkg.StructToJSONString(parserDSL)
	hash := pkg.HashSHA256(ts)
	key, ttl := http.GetIdempotencyKeyAndTTL(c)

	err := handler.Command.CreateOrCheckIdempotencyKey(ctx, organizationID, ledgerID, key, hash, ttl)
	if err != nil {
		mopentelemetry.HandleSpanError(&spanIdempotency, "Redis idempotency key", err)

		logger.Infof("Redis idempotency key: %v", err.Error())

		return http.WithError(c, err)
	}

	spanIdempotency.End()

	_, spanValidateDSL := tracer.Start(ctx, "handler.create_transaction_validate_dsl")

	validate, err := goldModel.ValidateSendSourceAndDistribute(parserDSL)
	if err != nil {
		mopentelemetry.HandleSpanError(&spanValidateDSL, "Failed to validate send source and distribute", err)

		logger.Error("Validation failed:", err.Error())

		return http.WithError(c, err)
	}

	spanValidateDSL.End()

	_, spanRaceCondition := tracer.Start(ctx, "handler.create_transaction.race_condition")

	handler.Command.AllKeysUnlocked(ctx, organizationID, ledgerID, validate.Aliases, hash)

	spanRaceCondition.End()

	ctxGetAccounts, spanGetAccounts := tracer.Start(ctx, "handler.create_transaction.get_accounts")

	token := http.GetTokenHeader(c)

	accounts, err := handler.getAccountsAndValidate(ctxGetAccounts, logger, token, hash, organizationID, ledgerID, validate, parserDSL)
	if err != nil {
		mopentelemetry.HandleSpanError(&spanGetAccounts, "Failed to get accounts", err)

		_, spanReleaseLock := tracer.Start(ctx, "handler.update_accounts.delete_locks_race_condition")
		handler.Command.DeleteLocks(ctx, organizationID, ledgerID, validate.Aliases, hash)

		spanReleaseLock.End()

		return http.WithError(c, err)
	}

	spanGetAccounts.End()

	ctxCreateTransaction, spanCreateTransaction := tracer.Start(ctx, "handler.create_transaction.create_transaction")

	err = mopentelemetry.SetSpanAttributesFromStruct(&spanCreateTransaction, "payload_command_create_transaction", parserDSL)
	if err != nil {
		mopentelemetry.HandleSpanError(&spanCreateTransaction, "Failed to convert parserDSL from struct to JSON string", err)

		return http.WithError(c, err)
	}

	tran, err := handler.Command.CreateTransaction(ctxCreateTransaction, organizationID, ledgerID, transactionID, &parserDSL)
	if err != nil {
		mopentelemetry.HandleSpanError(&spanCreateTransaction, "Failed to create transaction", err)

		logger.Error("Failed to create transaction", err.Error())

		_, spanReleaseLock := tracer.Start(ctx, "handler.update_accounts.delete_locks_race_condition")
		handler.Command.DeleteLocks(ctx, organizationID, ledgerID, validate.Aliases, hash)

		spanReleaseLock.End()

		return http.WithError(c, err)
	}

	spanCreateTransaction.End()

	e := make(chan error)
	result := make(chan []*operation.Operation)

	var operations []*operation.Operation

	ctxCreateOperation, spanCreateOperation := tracer.Start(ctx, "handler.create_transaction.create_operation")

	go handler.Command.CreateOperation(ctxCreateOperation, accounts, tran.ID, &parserDSL, *validate, result, e)
	select {
	case ops := <-result:
		operations = append(operations, ops...)
	case err := <-e:
		mopentelemetry.HandleSpanError(&spanCreateOperation, "Failed to create operations", err)

		logger.Error("Failed to create operations: ", err.Error())

		_, spanReleaseLock := tracer.Start(ctx, "handler.update_accounts.delete_locks_race_condition")
		handler.Command.DeleteLocks(ctx, organizationID, ledgerID, validate.Aliases, hash)

		spanReleaseLock.End()

		return http.WithError(c, err)
	}

	spanCreateOperation.End()

	ctxProcessAccounts, spanUpdateAccounts := tracer.Start(ctx, "handler.create_transaction.update_accounts")

	err = mopentelemetry.SetSpanAttributesFromStruct(&spanUpdateAccounts, "payload_handler_update_accounts", accounts)
	if err != nil {
		mopentelemetry.HandleSpanError(&spanUpdateAccounts, "Failed to convert accounts from struct to JSON string", err)
	}

	err = handler.Command.UpdateAccounts(ctxProcessAccounts, logger, *validate, token, organizationID, ledgerID, accounts)
	if err != nil {
		mopentelemetry.HandleSpanError(&spanUpdateAccounts, "Failed to update accounts", err)

		_, spanReleaseLock := tracer.Start(ctx, "handler.update_accounts.delete_locks_race_condition")
		handler.Command.DeleteLocks(ctx, organizationID, ledgerID, validate.Aliases, hash)

		spanReleaseLock.End()

		ctxUpdateTransactionStatus, spanUpdateTransactionStatus := tracer.Start(ctx, "handler.update_accounts.update_transaction_status")
		_, er := handler.Command.UpdateTransactionStatus(ctxUpdateTransactionStatus, organizationID, ledgerID, tran.IDtoUUID(), constant.DECLINED)

		if er != nil {
			mopentelemetry.HandleSpanError(&spanUpdateTransactionStatus, "Failed to update transaction status", err)

			logger.Errorf("Failed to update Transaction with ID: %s, Error: %s", tran.ID, err.Error())

			return http.WithError(c, er)
		}

		spanUpdateTransactionStatus.End()

		return http.WithError(c, err)
	}

	spanUpdateAccounts.End()

	ctxUpdateTransactionStatus, spanUpdateTransactionStatus := tracer.Start(ctx, "handler.create_transaction.update_transaction_status")
	_, err = handler.Command.UpdateTransactionStatus(ctxUpdateTransactionStatus, organizationID, ledgerID, tran.IDtoUUID(), constant.APPROVED)

	if err != nil {
		mopentelemetry.HandleSpanError(&spanUpdateTransactionStatus, "Failed to update transaction status", err)

		logger.Errorf("Failed to update Transaction with ID: %s, Error: %s", tran.ID, err.Error())

		_, spanReleaseLock := tracer.Start(ctx, "handler.update_accounts.delete_locks_race_condition")
		handler.Command.DeleteLocks(ctx, organizationID, ledgerID, validate.Aliases, hash)

		spanReleaseLock.End()

		return http.WithError(c, err)
	}

	spanUpdateTransactionStatus.End()

	ctxGetTransaction, spanGetTransaction := tracer.Start(ctx, "handler.create_transaction.get_transaction")

	tran, err = handler.Query.GetTransactionByID(ctxGetTransaction, organizationID, ledgerID, tran.IDtoUUID())
	if err != nil {
		mopentelemetry.HandleSpanError(&spanGetTransaction, "Failed to retrieve transaction", err)

		logger.Errorf("Failed to retrieve Transaction with ID: %s, Error: %s", tran.ID, err.Error())

		return http.WithError(c, err)
	}

	spanGetTransaction.End()

	tran.Source = validate.Sources
	tran.Destination = validate.Destinations
	tran.Operations = operations

	logger.Infof("Successfully updated Transaction with Organization ID: %s, Ledger ID: %s and ID: %s", organizationID.String(), ledgerID.String(), tran.ID)

	_, spanReleaseLock := tracer.Start(ctx, "handler.create_transaction.delete_race_condition")
	handler.Command.DeleteLocks(ctx, organizationID, ledgerID, validate.Aliases, hash)

	spanReleaseLock.End()

	go handler.logTransaction(ctx, operations, organizationID, ledgerID, tran.IDtoUUID())

	return http.Created(c, tran)
}

// getAccounts is a function that split aliases and ids, call the properly function and return Accounts
func (handler *TransactionHandler) getAccountsAndValidate(ctx context.Context, logger mlog.Logger, token, hash string, organizationID, ledgerID uuid.UUID, validate *goldModel.Responses, body goldModel.Transaction) ([]*account.Account, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	accounts, err := handler.Query.GetAccountsLedger(ctx, logger, token, organizationID, ledgerID, validate.Aliases)
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to get accounts", err)

		return nil, err
	}

	err = goldModel.ValidateAccounts(body, *validate, accounts)
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to validate accounts", err)

		handler.Command.DeleteLocks(ctx, organizationID, ledgerID, validate.Aliases, hash)

		return nil, err
	}

	searchAgain, err := handler.Command.LockBalanceVersion(ctx, organizationID, ledgerID, validate.Aliases, accounts)
	if err != nil {
		mopentelemetry.HandleSpanError(&span, "Failed to get account by alias gRPC on Ledger", err)

		logger.Error("Failed to get account by alias gRPC on Ledger", err.Error())

		return nil, err
	}

	if searchAgain {
		return handler.getAccountsAndValidate(ctx, logger, token, hash, organizationID, ledgerID, validate, body)
	}

	return accounts, nil
}

// logTransaction creates a message representing a transaction log and sends to auditing exchange
func (handler *TransactionHandler) logTransaction(ctx context.Context, operations []*operation.Operation, organizationID uuid.UUID, ledgerID uuid.UUID, transactionID uuid.UUID) {
	logger := pkg.NewLoggerFromContext(ctx)
	tracer := pkg.NewTracerFromContext(ctx)

	ctxLogTransaction, spanLogTransaction := tracer.Start(ctx, "handler.transaction.log_transaction")
	defer spanLogTransaction.End()

	if !isAuditLogEnabled() {
		logger.Infof("Audit logging not enabled. AUDIT_LOG_ENABLED='%s'", os.Getenv("AUDIT_LOG_ENABLED"))
		return
	}

	queueData := make([]mmodel.QueueData, 0)

	for _, o := range operations {
		oLog := o.ToLog()

		marshal, err := json.Marshal(oLog)
		if err != nil {
			logger.Fatalf("Failed to marshal operation to JSON string: %s", err.Error())
		}

		queueData = append(queueData, mmodel.QueueData{
			ID:    uuid.MustParse(o.ID),
			Value: marshal,
		})
	}

	queueMessage := mmodel.Queue{
		OrganizationID: organizationID,
		LedgerID:       ledgerID,
		AuditID:        transactionID,
		QueueData:      queueData,
	}

	if _, err := handler.Command.RabbitMQRepo.ProducerDefault(
		ctxLogTransaction,
		os.Getenv("RABBITMQ_AUDIT_EXCHANGE"),
		os.Getenv("RABBITMQ_AUDIT_KEY"),
		queueMessage,
	); err != nil {
		logger.Fatalf("Failed to send message: %s", err.Error())
	}
}

func isAuditLogEnabled() bool {
	envValue := strings.ToLower(strings.TrimSpace(os.Getenv("AUDIT_LOG_ENABLED")))
	return envValue != "false"
}
