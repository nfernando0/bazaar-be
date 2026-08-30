package config

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"bazar-be/models"
	"bazar-be/utils"

	"gorm.io/gorm"
)

// SeedAll populates all tables with realistic church / community bazaar dummy data
func SeedAll(db *gorm.DB, fresh bool) error {
	log.Println("[Seeder] Starting database seeding...")

	if fresh {
		log.Println("[Seeder] Resetting database tables (fresh mode)...")
		// Disable foreign key checks for clean truncation in MySQL
		db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
		db.Exec("TRUNCATE TABLE transaction_items;")
		db.Exec("TRUNCATE TABLE transactions;")
		db.Exec("TRUNCATE TABLE user_tokens;")
		db.Exec("TRUNCATE TABLE users;")
		db.Exec("TRUNCATE TABLE products;")
		db.Exec("TRUNCATE TABLE vendor_outlets;")
		db.Exec("TRUNCATE TABLE vendors;")
		db.Exec("TRUNCATE TABLE categories;")
		db.Exec("TRUNCATE TABLE outlets;")
		db.Exec("TRUNCATE TABLE bazaars;")
		db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
		log.Println("[Seeder] All tables truncated successfully.")
	}

	// 1. Seed Bazaars
	var bazaarCount int64
	db.Model(&models.Bazaar{}).Count(&bazaarCount)
	if bazaarCount == 0 || fresh {
		now := time.Now()
		bazaars := []models.Bazaar{
			{
				Name:        "Bazaar Akbar Paskah & UMKM Gereja 2026",
				Description: "Bazaar amal dan pameran aneka kuliner, souvenir, lilin doa, dan karya kasih paroki.",
				StartDate:   now.AddDate(0, -1, 0),
				EndDate:     now.AddDate(0, 1, 15),
				Status:      "active",
			},
			{
				Name:        "Bazaar Amal Kasih Natal 2025",
				Description: "Kegiatan penggalangan dana natal untuk beasiswa pendidikan anak panti asuhan.",
				StartDate:   now.AddDate(0, -4, 0),
				EndDate:     now.AddDate(0, -3, 0),
				Status:      "closed",
			},
			{
				Name:        "Bazaar Hari Pangan Sedunia 2026",
				Description: "Pameran hasil bumi, produk organik, dan olahan pangan sehat warga lingkungan.",
				StartDate:   now.AddDate(0, 2, 0),
				EndDate:     now.AddDate(0, 2, 14),
				Status:      "draft",
			},
		}
		if err := db.Create(&bazaars).Error; err != nil {
			return fmt.Errorf("failed to seed bazaars: %w", err)
		}
		log.Printf("[Seeder] Created %d bazaars.\n", len(bazaars))
	}

	// Retrieve primary active bazaar
	var activeBazaar models.Bazaar
	if err := db.Where("status = ?", "active").First(&activeBazaar).Error; err != nil {
		db.First(&activeBazaar)
	}

	// 2. Seed Outlets
	var outletCount int64
	db.Model(&models.Outlet{}).Count(&outletCount)
	if outletCount == 0 || fresh {
		outlets := []models.Outlet{
			{
				BazaarID: activeBazaar.ID,
				Name:     "Stan Makanan Utama & Nusantara",
				Code:     "OUT-FOOD-01",
				Location: "Lobby Utama Plaza Depan Gereja",
			},
			{
				BazaarID: activeBazaar.ID,
				Name:     "Stan Minuman & Cafe Senja",
				Code:     "OUT-BEV-02",
				Location: "Koridor Samping Aula Utama",
			},
			{
				BazaarID: activeBazaar.ID,
				Name:     "Stan Snack, Kue & Pastry",
				Code:     "OUT-SNACK-03",
				Location: "Area Tenda Sayap Barat",
			},
			{
				BazaarID: activeBazaar.ID,
				Name:     "Stan Lilin, Rohani & Souvenir",
				Code:     "OUT-SOUV-04",
				Location: "Galeri Rohani Gedung Karya Paroki",
			},
			{
				BazaarID: activeBazaar.ID,
				Name:     "Stan Pakaian, Craft & Thrift Amal",
				Code:     "OUT-CRAFT-05",
				Location: "Halaman Parkir Belakang",
			},
		}
		if err := db.Create(&outlets).Error; err != nil {
			return fmt.Errorf("failed to seed outlets: %w", err)
		}
		log.Printf("[Seeder] Created %d outlets.\n", len(outlets))
	}

	// Retrieve all outlets for foreign keys
	var allOutlets []models.Outlet
	db.Order("id ASC").Find(&allOutlets)

	// 3. Seed Categories
	var categoryCount int64
	db.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 || fresh {
		categories := []models.Category{
			{Name: "Makanan Berat & Nusantara"},
			{Name: "Minuman & Kopi"},
			{Name: "Snack & Jajanan Pasar"},
			{Name: "Pastry & Bakery"},
			{Name: "Barang Rohani & Lilin Doa"},
			{Name: "Merchandise & Souvenir"},
			{Name: "Pakaian & Kerajinan Tangan"},
		}
		if err := db.Create(&categories).Error; err != nil {
			return fmt.Errorf("failed to seed categories: %w", err)
		}
		log.Printf("[Seeder] Created %d categories.\n", len(categories))
	}

	var allCategories []models.Category
	db.Order("id ASC").Find(&allCategories)

	// 4. Seed Vendors
	var vendorCount int64
	db.Model(&models.Vendor{}).Count(&vendorCount)
	if vendorCount == 0 || fresh {
		vendors := []models.Vendor{
			{Name: "Dapur Mama Maria (Spesial Nasi & Gudeg)", ContactPerson: "Ibu Maria", Phone: "081234567801"},
			{Name: "Kopi Senja Express & Mocktail", ContactPerson: "Mas Andi", Phone: "081234567802"},
			{Name: "Pastry Berkah & Roti Jadul", ContactPerson: "Ibu Elisabeth", Phone: "081234567803"},
			{Name: "Kedai Dimsum & Siomay Oen", ContactPerson: "Koh Hendra", Phone: "081234567804"},
			{Name: "Kios Rohani Ave Maria & Lilin Doa", ContactPerson: "Sie Liturgi", Phone: "081234567805"},
			{Name: "Craft & Souvenir Karya Kasih OMK", ContactPerson: "Sie Kepemudaan", Phone: "081234567806"},
			{Name: "Sate & Ayam Bakar Pasundan", ContactPerson: "Pak Bambang", Phone: "081234567807"},
		}
		if err := db.Create(&vendors).Error; err != nil {
			return fmt.Errorf("failed to seed vendors: %w", err)
		}
		log.Printf("[Seeder] Created %d vendors.\n", len(vendors))
	}

	var allVendors []models.Vendor
	db.Order("id ASC").Find(&allVendors)

	// 5. Seed VendorOutlets (Mapping)
	var voCount int64
	db.Model(&models.VendorOutlet{}).Count(&voCount)
	if voCount == 0 || fresh {
		if len(allVendors) >= 7 && len(allOutlets) >= 5 {
			vendorOutlets := []models.VendorOutlet{
				{VendorID: allVendors[0].ID, OutletID: allOutlets[0].ID, BoothNumber: "Booth-A1"},
				{VendorID: allVendors[6].ID, OutletID: allOutlets[0].ID, BoothNumber: "Booth-A2"},
				{VendorID: allVendors[1].ID, OutletID: allOutlets[1].ID, BoothNumber: "Booth-B1"},
				{VendorID: allVendors[3].ID, OutletID: allOutlets[1].ID, BoothNumber: "Booth-B2"},
				{VendorID: allVendors[2].ID, OutletID: allOutlets[2].ID, BoothNumber: "Booth-C1"},
				{VendorID: allVendors[4].ID, OutletID: allOutlets[3].ID, BoothNumber: "Booth-D1"},
				{VendorID: allVendors[5].ID, OutletID: allOutlets[4].ID, BoothNumber: "Booth-E1"},
			}
			if err := db.Create(&vendorOutlets).Error; err != nil {
				return fmt.Errorf("failed to seed vendor outlets: %w", err)
			}
			log.Printf("[Seeder] Created %d vendor outlet assignments.\n", len(vendorOutlets))
		}
	}

	// 6. Seed Products
	var productCount int64
	db.Model(&models.Product{}).Count(&productCount)
	if productCount == 0 || fresh {
		if len(allVendors) >= 7 && len(allCategories) >= 7 {
			cFood := allCategories[0].ID
			cBev := allCategories[1].ID
			cSnack := allCategories[2].ID
			cBakery := allCategories[3].ID
			cRel := allCategories[4].ID
			cSouv := allCategories[5].ID
			cCloth := allCategories[6].ID

			products := []models.Product{
				// Vendor 0: Dapur Mama Maria
				{VendorID: allVendors[0].ID, CategoryID: &cFood, Name: "Nasi Campur Bali Komplit", Price: 35000, Stock: 45},
				{VendorID: allVendors[0].ID, CategoryID: &cFood, Name: "Nasi Gudeg Telur Tahu Krecek", Price: 28000, Stock: 50},
				{VendorID: allVendors[0].ID, CategoryID: &cFood, Name: "Nasi Ayam Bakar Madu", Price: 30000, Stock: 60},

				// Vendor 6: Sate & Ayam Bakar Pasundan
				{VendorID: allVendors[6].ID, CategoryID: &cFood, Name: "Sate Ayam Bumbu Kacang (10 Tusuk)", Price: 25000, Stock: 70},
				{VendorID: allVendors[6].ID, CategoryID: &cFood, Name: "Tahu Tempe Goreng Lengkuas", Price: 10000, Stock: 80},

				// Vendor 1: Kopi Senja Express
				{VendorID: allVendors[1].ID, CategoryID: &cBev, Name: "Es Kopi Susu Gula Aren Senja", Price: 18000, Stock: 100},
				{VendorID: allVendors[1].ID, CategoryID: &cBev, Name: "Iced Matcha Green Tea Latte", Price: 22000, Stock: 80},
				{VendorID: allVendors[1].ID, CategoryID: &cBev, Name: "Lemon Tea Madu Segar", Price: 12000, Stock: 120},
				{VendorID: allVendors[1].ID, CategoryID: &cBev, Name: "Es Cokelat Klasik Creamy", Price: 15000, Stock: 90},

				// Vendor 3: Kedai Dimsum & Siomay Oen
				{VendorID: allVendors[3].ID, CategoryID: &cSnack, Name: "Dimsum Ayam Udang Komplit (4 Pcs)", Price: 22000, Stock: 80},
				{VendorID: allVendors[3].ID, CategoryID: &cSnack, Name: "Siomay Bandung Ikan Tenggiri", Price: 20000, Stock: 50},

				// Vendor 2: Pastry Berkah & Roti Jadul
				{VendorID: allVendors[2].ID, CategoryID: &cBakery, Name: "Roti Sisir Mentega Jadul", Price: 15000, Stock: 40},
				{VendorID: allVendors[2].ID, CategoryID: &cBakery, Name: "Croissant Butter Wangi", Price: 20000, Stock: 35},
				{VendorID: allVendors[2].ID, CategoryID: &cSnack, Name: "Risoles Rogout Ayam Sayur (2 Pcs)", Price: 12000, Stock: 60},
				{VendorID: allVendors[2].ID, CategoryID: &cSnack, Name: "Pastel Panggang Telur Daging", Price: 12000, Stock: 50},

				// Vendor 4: Kios Rohani Ave Maria
				{VendorID: allVendors[4].ID, CategoryID: &cRel, Name: "Lilin Doa Kaca Merah Besar", Price: 25000, Stock: 80},
				{VendorID: allVendors[4].ID, CategoryID: &cRel, Name: "Lilin Doa Novena Putih", Price: 15000, Stock: 100},
				{VendorID: allVendors[4].ID, CategoryID: &cRel, Name: "Rosario Kayu Zaitun Asli", Price: 35000, Stock: 40},
				{VendorID: allVendors[4].ID, CategoryID: &cRel, Name: "Patung Mini Keluarga Kudus Resin", Price: 45000, Stock: 30},

				// Vendor 5: Craft & Souvenir Karya Kasih
				{VendorID: allVendors[5].ID, CategoryID: &cCloth, Name: "Kaos Polo Panitia Bazar Paskah", Price: 85000, Stock: 50},
				{VendorID: allVendors[5].ID, CategoryID: &cSouv, Name: "Tote Bag Kanvas Motif Kasih", Price: 35000, Stock: 60},
				{VendorID: allVendors[5].ID, CategoryID: &cSouv, Name: "Gantungan Kunci Rohani Kayu", Price: 15000, Stock: 100},
			}

			if err := db.Create(&products).Error; err != nil {
				return fmt.Errorf("failed to seed products: %w", err)
			}
			log.Printf("[Seeder] Created %d products.\n", len(products))
		}
	}

	var allProducts []models.Product
	db.Order("id ASC").Find(&allProducts)

	// 7. Seed Users (Admin & Multi-Outlet Cashiers)
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 || fresh {
		defaultPassword, _ := utils.HashPassword("password123")

		var o1, o2, o3, o4, o5 *uint
		if len(allOutlets) >= 5 {
			o1 = &allOutlets[0].ID
			o2 = &allOutlets[1].ID
			o3 = &allOutlets[2].ID
			o4 = &allOutlets[3].ID
			o5 = &allOutlets[4].ID
		}

		users := []models.User{
			// Admins
			{
				Name:         "Nando Admin Master",
				Email:        "admin@example.com",
				PasswordHash: defaultPassword,
				Role:         "admin",
				OutletID:     nil,
			},
			{
				Name:         "Koordinator Sie Acara",
				Email:        "panitia@bazar.com",
				PasswordHash: defaultPassword,
				Role:         "admin",
				OutletID:     nil,
			},
			// Cashiers
			{
				Name:         "Kasir Stan Makanan",
				Email:        "kasir@example.com",
				PasswordHash: defaultPassword,
				Role:         "cashier",
				OutletID:     o1,
			},
			{
				Name:         "Kasir Stan Makanan 1",
				Email:        "kasir1@example.com",
				PasswordHash: defaultPassword,
				Role:         "cashier",
				OutletID:     o1,
			},
			{
				Name:         "Kasir Stan Minuman",
				Email:        "kasir2@example.com",
				PasswordHash: defaultPassword,
				Role:         "cashier",
				OutletID:     o2,
			},
			{
				Name:         "Kasir Stan Snack & Kue",
				Email:        "kasir3@example.com",
				PasswordHash: defaultPassword,
				Role:         "cashier",
				OutletID:     o3,
			},
			{
				Name:         "Kasir Stan Rohani & Souvenir",
				Email:        "kasir4@example.com",
				PasswordHash: defaultPassword,
				Role:         "cashier",
				OutletID:     o4,
			},
			{
				Name:         "Kasir Stan Pakaian & Craft",
				Email:        "kasir5@example.com",
				PasswordHash: defaultPassword,
				Role:         "cashier",
				OutletID:     o5,
			},
		}

		if err := db.Create(&users).Error; err != nil {
			return fmt.Errorf("failed to seed users: %w", err)
		}
		log.Printf("[Seeder] Created %d users.\n", len(users))
	}

	var allUsers []models.User
	db.Order("id ASC").Find(&allUsers)

	// 8. Seed Realistic Transactions & Transaction Items
	var txCount int64
	db.Model(&models.Transaction{}).Count(&txCount)
	if txCount == 0 || fresh {
		if len(allOutlets) >= 5 && len(allProducts) >= 15 && len(allUsers) >= 3 {
			log.Println("[Seeder] Seeding realistic transactions and transaction items...")

			// Map products by outlet ID for coherent sales records
			outletProductMap := map[uint][]models.Product{
				allOutlets[0].ID: {allProducts[0], allProducts[1], allProducts[2], allProducts[3], allProducts[4]},
				allOutlets[1].ID: {allProducts[5], allProducts[6], allProducts[7], allProducts[8], allProducts[9], allProducts[10]},
				allOutlets[2].ID: {allProducts[11], allProducts[12], allProducts[13], allProducts[14]},
				allOutlets[3].ID: {allProducts[15], allProducts[16], allProducts[17], allProducts[18]},
				allOutlets[4].ID: {allProducts[19], allProducts[20], allProducts[21]},
			}

			// Map cashiers to outlets
			outletCashierMap := map[uint]uint{}
			for _, u := range allUsers {
				if u.OutletID != nil {
					outletCashierMap[*u.OutletID] = u.ID
				}
			}

			paymentMethods := []string{"cash", "qris", "midtrans", "cash", "qris"}
			now := time.Now()

			// Seed 30 transactions spanning the last 5 days
			trxIndex := 100
			for dayOffset := 4; dayOffset >= 0; dayOffset-- {
				trxDate := now.AddDate(0, 0, -dayOffset)

				for _, outlet := range allOutlets {
					prods := outletProductMap[outlet.ID]
					if len(prods) == 0 {
						continue
					}

					cashierID := outletCashierMap[outlet.ID]
					if cashierID == 0 {
						cashierID = allUsers[0].ID
					}

					// 1 to 3 transactions per outlet per day
					numTrx := rand.Intn(2) + 1
					for i := 0; i < numTrx; i++ {
						trxIndex++
						trxHour := 9 + rand.Intn(10)
						trxMinute := rand.Intn(60)
						tTime := time.Date(trxDate.Year(), trxDate.Month(), trxDate.Day(), trxHour, trxMinute, rand.Intn(60), 0, time.Local)

						payMethod := paymentMethods[rand.Intn(len(paymentMethods))]
						trxCode := fmt.Sprintf("TRX-%s-%04d", tTime.Format("20060102-150405"), trxIndex)

						// 1 to 3 distinct items in transaction
						numItems := rand.Intn(2) + 1
						var totalAmount float64
						var items []models.TransactionItem

						perm := rand.Perm(len(prods))
						for j := 0; j < numItems && j < len(prods); j++ {
							prod := prods[perm[j]]
							qty := rand.Intn(3) + 1
							subtotal := float64(qty) * prod.Price
							totalAmount += subtotal

							items = append(items, models.TransactionItem{
								ProductID:   prod.ID,
								ProductName: prod.Name,
								Quantity:    qty,
								UnitPrice:   prod.Price,
								Subtotal:    subtotal,
								CreatedAt:   tTime,
								UpdatedAt:   tTime,
							})
						}

						trx := models.Transaction{
							OutletID:        outlet.ID,
							CashierID:       cashierID,
							TransactionCode: trxCode,
							TotalAmount:     totalAmount,
							PaymentMethod:   payMethod,
							Items:           items,
							CreatedAt:       tTime,
							UpdatedAt:       tTime,
						}

						if err := db.Create(&trx).Error; err != nil {
							log.Printf("[Seeder] Warning: Failed to seed transaction %s: %v", trxCode, err)
						}
					}
				}
			}

			var finalTrxCount int64
			db.Model(&models.Transaction{}).Count(&finalTrxCount)
			log.Printf("[Seeder] Successfully seeded %d transactions.\n", finalTrxCount)
		}
	}

	log.Println("[Seeder] Database seeding completed successfully! 🎉")
	return nil
}
