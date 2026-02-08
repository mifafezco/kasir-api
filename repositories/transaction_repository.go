package repositories

import (
    "database/sql"
    "fmt" // Ditambahkan karena ada fmt.Errorf
    "kasir-api/model" // Pastikan konsisten menggunakan 'model'
)

type TransactionRepository struct {
    db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
    return &TransactionRepository{db: db}
}

// Gunakan model. (sesuai nama package di import)
func (repo *TransactionRepository) CreateTransaction(items []model.CheckoutItem) (*model.Transaction, error) {
    tx, err := repo.db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    totalAmount := 0
    details := make([]model.TransactionDetail, 0)

    for _, item := range items {
        var productPrice, stock int
        var productName string

        err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.CategoriesID).Scan(&productName, &productPrice, &stock)
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("product id %d not found", item.CategoriesID)
        }
        if err != nil {
            return nil, err
        }

        // Cek stok sebelum transaksi (Saran Tambahan)
        if stock < item.Quantity {
            return nil, fmt.Errorf("insufficient stock for product %s", productName)
        }
  
        subtotal := productPrice * item.Quantity
        totalAmount += subtotal

        _, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.CategoriesID)
        if err != nil {
            return nil, err
        }

        details = append(details, model.TransactionDetail{
            // Pastikan field CategoriesID ada di struct model.TransactionDetail
            // Jika di model tadi namanya CategoriesID, sesuaikan di sini
            CategoriesID:   item.CategoriesID, 
            ProductName: productName,
            Quantity:    item.Quantity,
            Subtotal:    subtotal,
        })
    }

    var transactionID int
    err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
    if err != nil {
        return nil, err
    }


    for i := range details {
    details[i].TransactionID = transactionID
    
    // Pastikan nama kolom di DB (product_id vs categories_id) konsisten dengan model kamu
    // Jika di struct model.TransactionDetail namanya CategoriesID, maka gunakan itu:
    _, err = tx.Exec(`
        INSERT INTO transaction_details (transaction_id, categories_id, quantity, subtotal) 
        VALUES ($1, $2, $3, $4)`,
        transactionID, 
        details[i].CategoriesID, // Sesuaikan dengan nama field di struct model
        details[i].Quantity, 
        details[i].Subtotal,
    )
    
    if err != nil {
        // Jika gagal di sini, tx.Rollback() di defer akan membatalkan semuanya
        return nil, fmt.Errorf("failed to insert detail for category %d: %v", details[i].CategoriesID, err)
    }
    }

    if err := tx.Commit(); err != nil {
        return nil, err
    }

    return &model.Transaction{
        ID:          transactionID,
        TotalAmount: totalAmount,
        Details:     details,
    }, nil
}

func (repo *TransactionRepository) GetReport(startDate, endDate string) (*model.ReportResponse, error) {
	var report model.ReportResponse

	// 1. Hitung Total Revenue dan Total Transaksi
	queryStats := `
		SELECT COALESCE(SUM(total_amount), 0), COUNT(id) 
		FROM transactions 
		WHERE created_at::date BETWEEN $1 AND $2`
	
	err := repo.db.QueryRow(queryStats, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// 2. Cari Kategori Terlaris (Join dengan tabel categories)
	queryBest := `
		SELECT c.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN categories c ON td.categories_id = c.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at::date BETWEEN $1 AND $2
		GROUP BY c.name
		ORDER BY total_qty DESC
		LIMIT 1`

	err = repo.db.QueryRow(queryBest, startDate, endDate).Scan(
		&report.BestCategory.Nama, 
		&report.BestCategory.QtyTerjual,
	)
	
	// Jika belum ada transaksi, handle agar tidak error scan
	if err == sql.ErrNoRows {
		report.BestCategory = model.BestCategory{Nama: "-", QtyTerjual: 0}
	} else if err != nil {
		return nil, err
	}

	return &report, nil
}