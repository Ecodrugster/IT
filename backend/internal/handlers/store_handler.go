package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

// StoreItem represents a product on the college shop
type StoreItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CoinsPrice  int    `json:"coins_price"`
	StarsPrice  int    `json:"stars_price"`
	Image       string `json:"image"` // Emoji representing the item
	Stock       int    `json:"stock"`
}

// Purchase represents a logged item purchase in Firestore
type Purchase struct {
	ID         string    `json:"id"`
	StudentID  string    `json:"student_id"`
	ItemID     string    `json:"item_id"`
	ItemName   string    `json:"item_name"`
	ItemIcon   string    `json:"item_icon"`
	CoinsSpent int       `json:"coins_spent"`
	StarsSpent int       `json:"stars_spent"`
	PromoCode  string    `json:"promo_code"`
	Status     string    `json:"status"` // "active" | "claimed"
	CreatedAt  time.Time `json:"created_at"`
}

// AvailableStoreItems contains the list of catalog items
var AvailableStoreItems = []StoreItem{
	{ID: "hoodie", Name: "Худи ITSTEP Premium", Description: "Премиальное теплое худи с вышитым фирменным логотипом ITSTEP", CoinsPrice: 400, StarsPrice: 3, Image: "/images/store/hoodie.png", Stock: 15},
	{ID: "usb", Name: "Флешка ITSTEP 64GB", Description: "Быстрая флешка USB 3.0 на 64 ГБ для ваших курсовых и проектов", CoinsPrice: 200, StarsPrice: 0, Image: "/images/store/usb.png", Stock: 50},
	{ID: "notebook", Name: "Фирменная тетрадь", Description: "Стильная тетрадь для конспектов лекций в клетку", CoinsPrice: 100, StarsPrice: 0, Image: "/images/store/notepad.png", Stock: 100},
	{ID: "tshirt", Name: "Футболка ITSTEP Tech", Description: "Хлопковая дышащая футболка с принтом ITSTEP Tech", CoinsPrice: 150, StarsPrice: 1, Image: "👕", Stock: 30},
	{ID: "thermos", Name: "Термокружка ITSTEP", Description: "Сохранит ваш кофе или чай горячим на всех парах", CoinsPrice: 200, StarsPrice: 2, Image: "☕", Stock: 20},
	{ID: "stickers", Name: "Набор стикеров ITSTEP", Description: "Яркие стикеры для кастомизации ноутбука, планшета или телефона", CoinsPrice: 20, StarsPrice: 0, Image: "🏷️", Stock: 200},
	{ID: "backpack", Name: "Рюкзак ITSTEP Urban", Description: "Премиальный городской рюкзак для ноутбука с влагозащитой", CoinsPrice: 500, StarsPrice: 5, Image: "/images/store/backpack.png", Stock: 10},
	{ID: "mouse", Name: "Беспроводная мышь", Description: "Эргономичная беспроводная мышь с фирменной гравировкой ITSTEP", CoinsPrice: 250, StarsPrice: 1, Image: "🖱️", Stock: 25},
	{ID: "keyboard", Name: "Механическая клавиатура", Description: "Игровая механическая клавиатура с RGB-подсветкой и тихим откликом", CoinsPrice: 800, StarsPrice: 8, Image: "/images/store/klava.png", Stock: 5},
	{ID: "cap", Name: "Кепка ITSTEP Cap", Description: "Стильная хлопковая кепка с вышитым логотипом ITSTEP", CoinsPrice: 180, StarsPrice: 1, Image: "/images/store/cap.png", Stock: 40},
	{ID: "bottle", Name: "Спортивная бутылка", Description: "Алюминиевая спортивная бутылка для воды на 750 мл", CoinsPrice: 120, StarsPrice: 0, Image: "/images/store/butilka.png", Stock: 35},
	{ID: "notepad", Name: "Кожаный ежедневник", Description: "Деловой кожаный блокнот-ежедневник в комплекте с ручкой", CoinsPrice: 150, StarsPrice: 0, Image: "/images/store/notepad.png", Stock: 60},
}

// GetStoreItems lists all products in the catalog
func GetStoreItems(c *gin.Context) {
	c.JSON(http.StatusOK, AvailableStoreItems)
}

// PurchaseItem processes a student purchase
func PurchaseItem(c *gin.Context) {
	uid, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		ItemID string `json:"item_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан ID товара"})
		return
	}

	// Find the item
	var targetItem *StoreItem
	for i := range AvailableStoreItems {
		if AvailableStoreItems[i].ID == req.ItemID {
			targetItem = &AvailableStoreItems[i]
			break
		}
	}
	if targetItem == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
		return
	}

	// Fetch student document
	studentRef := repositories.FirestoreClient.Collection("users").Doc(uid)
	studentDoc, err := studentRef.Get(c.Request.Context())
	if err != nil || !studentDoc.Exists() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Профиль студента не найден"})
		return
	}

	data := studentDoc.Data()
	var coins, stars int
	if cVal, ok := data["coins"].(int64); ok {
		coins = int(cVal)
	} else if cValFloat, ok := data["coins"].(float64); ok {
		coins = int(cValFloat)
	}
	if sVal, ok := data["stars"].(int64); ok {
		stars = int(sVal)
	} else if sValFloat, ok := data["stars"].(float64); ok {
		stars = int(sValFloat)
	}

	// Validate balances
	if coins < targetItem.CoinsPrice || stars < targetItem.StarsPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недостаточно монет 🪙 или звезд 🌟 для покупки этого товара"})
		return
	}

	// Subtract price
	coins -= targetItem.CoinsPrice
	stars -= targetItem.StarsPrice

	// Generate a secure unique promo code
	promoBytes := make([]byte, 3)
	_, _ = rand.Read(promoBytes)
	promoCode := fmt.Sprintf("STEP-%X", promoBytes)

	// Save student balance updates
	_, err = studentRef.Set(c.Request.Context(), map[string]interface{}{
		"coins": coins,
		"stars": stars,
	}, firestore.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось списать монеты/звезды"})
		return
	}

	// Save purchase record in Firestore
	purchaseData := map[string]interface{}{
		"student_id":  uid,
		"item_id":     targetItem.ID,
		"item_name":   targetItem.Name,
		"item_icon":   targetItem.Image,
		"coins_spent": targetItem.CoinsPrice,
		"stars_spent": targetItem.StarsPrice,
		"promo_code":  promoCode,
		"status":      "active",
		"created_at":  time.Now(),
	}
	
	_, _, err = repositories.FirestoreClient.Collection("purchases").Add(c.Request.Context(), purchaseData)
	if err != nil {
		// Log the error but student already lost coins, we must return the code
		fmt.Printf("[Store] Failed to save purchase doc: %v\n", err)
	}

	// Send purchase notification
	notifTitle := "Успешная покупка! 🛍️"
	notifMsg := fmt.Sprintf("Вы приобрели «%s» за %d монет и %d звезд. Ваш промокод: %s. Покажите его в деканате!", targetItem.Name, targetItem.CoinsPrice, targetItem.StarsPrice, promoCode)
	_ = CreateNotification(c.Request.Context(), uid, "system", notifTitle, notifMsg, "/store")

	c.JSON(http.StatusOK, gin.H{
		"message":     "Покупка оформлена успешно!",
		"promo_code":  promoCode,
		"coins_left":  coins,
		"stars_left":  stars,
		"item_name":   targetItem.Name,
		"item_icon":   targetItem.Image,
	})
}

// GetPurchaseHistory lists all purchases made by the current student
func GetPurchaseHistory(c *gin.Context) {
	uid, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	iter := repositories.FirestoreClient.Collection("purchases").
		Where("student_id", "==", uid).
		Documents(c.Request.Context())

	var history []Purchase
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось загрузить историю покупок"})
			return
		}

		data := doc.Data()
		var createdAt time.Time
		if t, ok := data["created_at"].(time.Time); ok {
			createdAt = t
		}

		p := Purchase{
			ID:         doc.Ref.ID,
			StudentID:  asString(data["student_id"]),
			ItemID:     asString(data["item_id"]),
			ItemName:   asString(data["item_name"]),
			ItemIcon:   asString(data["item_icon"]),
			CoinsSpent: asInt(data["coins_spent"]),
			StarsSpent: asInt(data["stars_spent"]),
			PromoCode:  asString(data["promo_code"]),
			Status:     asString(data["status"]),
			CreatedAt:  createdAt,
		}
		history = append(history, p)
	}

	// Sort by Date descending
	sort.Slice(history, func(i, j int) bool {
		return history[i].CreatedAt.After(history[j].CreatedAt)
	})

	c.JSON(http.StatusOK, history)
}

// GetAdminPurchases lists all student purchases (Admin only)
func GetAdminPurchases(c *gin.Context) {
	iter := repositories.FirestoreClient.Collection("purchases").Documents(c.Request.Context())
	var results []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch purchases"})
			return
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID

		studentID := asString(data["student_id"])
		studentDoc, err := repositories.FirestoreClient.Collection("users").Doc(studentID).Get(c.Request.Context())
		if err == nil && studentDoc.Exists() {
			sData := studentDoc.Data()
			data["student_name"] = getDisplayNameFromUserData(sData, studentID)
			data["student_email"] = asString(sData["email"])
			data["student_group"] = asString(sData["group_name"])
			if data["student_group"] == "" {
				data["student_group"] = asString(sData["group"])
			}
		} else {
			data["student_name"] = "Студент (" + studentID + ")"
		}

		results = append(results, data)
	}

	sort.Slice(results, func(i, j int) bool {
		tI, okI := results[i]["created_at"].(time.Time)
		tJ, okJ := results[j]["created_at"].(time.Time)
		if !okI || !okJ {
			return false
		}
		return tI.After(tJ)
	})

	c.JSON(http.StatusOK, results)
}

// ClaimPurchase marks a promo code voucher as claimed and delivers the item
func ClaimPurchase(c *gin.Context) {
	var req struct {
		PromoCode  string `json:"promo_code"`
		PurchaseID string `json:"purchase_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var ref *firestore.DocumentRef
	if req.PurchaseID != "" {
		ref = repositories.FirestoreClient.Collection("purchases").Doc(req.PurchaseID)
	} else if req.PromoCode != "" {
		iter := repositories.FirestoreClient.Collection("purchases").
			Where("promo_code", "==", strings.ToUpper(strings.TrimSpace(req.PromoCode))).
			Limit(1).
			Documents(c.Request.Context())
		doc, err := iter.Next()
		if err == iterator.Done {
			c.JSON(http.StatusNotFound, gin.H{"error": "Промокод не найден"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска промокода"})
			return
		}
		ref = doc.Ref
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Укажите промокод или ID покупки"})
		return
	}

	doc, err := ref.Get(c.Request.Context())
	if err != nil || !doc.Exists() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Покупка не найдена"})
		return
	}

	data := doc.Data()
	if asString(data["status"]) == "claimed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Этот промокод уже был использован и товар выдан!"})
		return
	}

	_, err = ref.Update(c.Request.Context(), []firestore.Update{
		{Path: "status", Value: "claimed"},
		{Path: "claimed_at", Value: time.Now()},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить статус"})
		return
	}

	studentID := asString(data["student_id"])
	itemName := asString(data["item_name"])

	notifTitle := "Товар успешно выдан! 🎒"
	notifMsg := fmt.Sprintf("Ваш промокод %s на «%s» был успешно активирован в деканате. Товар выдан! Спасибо за вашу активность!", asString(data["promo_code"]), itemName)
	_ = CreateNotification(c.Request.Context(), studentID, "system", notifTitle, notifMsg, "/store")

	c.JSON(http.StatusOK, gin.H{
		"message":    "Товар успешно выдан студенту!",
		"item_name":  itemName,
		"promo_code": asString(data["promo_code"]),
	})
}
