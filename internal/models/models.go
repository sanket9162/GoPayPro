package models

import (
	"context"
	"database/sql"
	"time"
)

// DBModel is the type for database connection values
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for all models
type Models struct {
	DB DBModel
}

// NewModels retruns a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Widget is the type for all widgets
type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	IsRecurring    bool      `json:"is_recurring"`
	PlanID         string    `json:"plan_id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// Order is the type for all orders
type Order struct {
	ID            int       `json:"id"`
	WidgetID      int       `json:"widget_id"`
	TransactionID int       `json:"transaction_id"`
	CustomerID    int       `json:"customer_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreateAt      time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

// Status is the type for statues
type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreateAt  time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TransactionStatus is the type for transaction statues
type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreateAt  time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Transaction is the type for transaction
type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
	BankReturnCode      string    `json:"Bank_return_code"`
	TransactionStatusID int       `json:"Transaction_status_id"`
	CreateAt            time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

// User is the type for User
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"First_name"`
	LastName  string    `json:"Last_Name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreateAt  time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"First_name"`
	LastName  string    `json:"Last_Name"`
	Email     string    `json:"email"`
	CreateAt  time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widget
	row := m.DB.QueryRowContext(ctx, `
	select 
		id, name, description, inventory_level, price, coalesce(image, ''),
		is_recurring, plan_id,
		created_at, updated_at
	from 
		widgets
	where id = ?`, id)
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.IsRecurring,
		&widget.PlanID,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)
	if err != nil {
		return widget, err
	}

	return widget, nil
}

// InsertTransaction inserts a new txn, and returns its id
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into transactions
			(amount, currency, last_four, bank_return_code, expiry_month,  expiry_year,
			payment_intent, payment_method,
			transaction_status_id, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := m.DB.ExecContext(ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.PaymentIntent,
		txn.PaymentMethod,
		txn.TransactionStatusID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertOrder inserts a new order, and returns its id
func (m *DBModel) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into orders
			(widget_id, transaction_id, status_id, quantity, customer_id,
			amount, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := m.DB.ExecContext(ctx, stmt,
		order.WidgetID,
		order.TransactionID,
		order.StatusID,
		order.CustomerID,
		order.Quantity,
		order.Amount,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertCustomer inserts a customer, and returns its id
func (m *DBModel) InsertCustomer(c Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into customers
			(first_name, last_name, email, created_at, updated_at)
		values (?, ?, ?, ?, ?)
	`
	result, err := m.DB.ExecContext(ctx, stmt,
		c.FirstName,
		c.LastName,
		c.Email,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
