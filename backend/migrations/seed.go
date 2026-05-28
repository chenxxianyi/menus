package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"menu-recommend/config"
	"menu-recommend/internal/model"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Seed categories
	categories := []model.RecipeCategory{
		{Name: "家常菜", Icon: "home-cooking", Sort: 1},
		{Name: "快手菜", Icon: "quick", Sort: 2},
		{Name: "早餐", Icon: "breakfast", Sort: 3},
		{Name: "汤羹", Icon: "soup", Sort: 4},
		{Name: "主食", Icon: "staple", Sort: 5},
		{Name: "凉菜", Icon: "cold-dish", Sort: 6},
		{Name: "烘焙", Icon: "baking", Sort: 7},
		{Name: "饮品", Icon: "drink", Sort: 8},
	}
	for _, c := range categories {
		db.Where("name = ?", c.Name).FirstOrCreate(&c)
	}

	// Seed recipes
	recipes := []model.Recipe{
		{
			Title:       "番茄炒蛋",
			Description: "经典家常菜，简单快手，营养美味",
			CategoryID:  1,
			CookTime:    15,
			Difficulty:  "简单",
			PeopleCount: 2,
			Taste:       "咸鲜",
			HealthTags:  model.JSON(`["快手","低脂"]`),
			Ingredients: model.JSON(`[{"name":"鸡蛋","amount":"3个"},{"name":"番茄","amount":"2个"}]`),
			Seasonings:  model.JSON(`[{"name":"盐","amount":"适量"},{"name":"糖","amount":"少许"},{"name":"葱","amount":"适量"}]`),
			Steps:       model.JSON(`[{"step":1,"desc":"鸡蛋打散，番茄切块"},{"step":2,"desc":"热锅凉油，倒入蛋液炒熟盛出"},{"step":3,"desc":"锅中留底油，放入番茄翻炒出汁"},{"step":4,"desc":"加入鸡蛋，调入盐和糖，翻炒均匀"}]`),
			Tips:        "番茄用开水烫一下更容易去皮",
			Nutrition:   model.JSON(`{"calories":150,"protein":12,"fat":8,"carbs":10}`),
		},
		{
			Title:       "红烧排骨",
			Description: "经典家常菜，软烂入味，色泽红亮",
			CategoryID:  1,
			CookTime:    40,
			Difficulty:  "中等",
			PeopleCount: 3,
			Taste:       "咸鲜",
			HealthTags:  model.JSON(`["高蛋白"]`),
			Ingredients: model.JSON(`[{"name":"排骨","amount":"500g"},{"name":"冰糖","amount":"30g"},{"name":"姜片","amount":"3片"}]`),
			Seasonings:  model.JSON(`[{"name":"生抽","amount":"2勺"},{"name":"老抽","amount":"1勺"},{"name":"料酒","amount":"1勺"}]`),
			Steps:       model.JSON(`[{"step":1,"desc":"排骨冷水下锅焯水去腥"},{"step":2,"desc":"锅中放油，加冰糖炒糖色"},{"step":3,"desc":"下排骨翻炒上色"},{"step":4,"desc":"加调料和热水，小火炖30分钟"},{"step":5,"desc":"大火收汁即可"}]`),
			Tips:        "一定要加热水，冷水会让肉质收紧",
			Nutrition:   model.JSON(`{"calories":520,"protein":35,"fat":28,"carbs":42}`),
		},
		{
			Title:       "清炒时蔬",
			Description: "清淡爽口，营养丰富",
			CategoryID:  2,
			CookTime:    10,
			Difficulty:  "简单",
			PeopleCount: 2,
			Taste:       "清淡",
			HealthTags:  model.JSON(`["低脂","素食","快手"]`),
			Ingredients: model.JSON(`[{"name":"西兰花","amount":"1颗"},{"name":"胡萝卜","amount":"1根"}]`),
			Seasonings:  model.JSON(`[{"name":"盐","amount":"适量"},{"name":"蒜","amount":"3瓣"}]`),
			Steps:       model.JSON(`[{"step":1,"desc":"蔬菜洗净切好"},{"step":2,"desc":"热锅凉油，爆香蒜末"},{"step":3,"desc":"下蔬菜翻炒至断生"},{"step":4,"desc":"调入盐，翻炒均匀出锅"}]`),
			Nutrition:   model.JSON(`{"calories":80,"protein":4,"fat":3,"carbs":12}`),
		},
		{
			Title:       "糖醋里脊",
			Description: "酸甜可口，外酥里嫩",
			CategoryID:  1,
			CookTime:    30,
			Difficulty:  "中等",
			PeopleCount: 3,
			Taste:       "酸甜",
			HealthTags:  model.JSON(`["高蛋白"]`),
			Ingredients: model.JSON(`[{"name":"猪里脊","amount":"300g"},{"name":"淀粉","amount":"适量"},{"name":"鸡蛋","amount":"1个"}]`),
			Seasonings:  model.JSON(`[{"name":"醋","amount":"3勺"},{"name":"糖","amount":"3勺"},{"name":"番茄酱","amount":"2勺"}]`),
			Steps:       model.JSON(`[{"step":1,"desc":"里脊切条，用盐、料酒腌制15分钟"},{"step":2,"desc":"裹蛋液和淀粉，炸至金黄"},{"step":3,"desc":"调糖醋汁"},{"step":4,"desc":"锅中倒入糖醋汁，下里脊翻炒均匀"}]`),
			Nutrition:   model.JSON(`{"calories":380,"protein":25,"fat":18,"carbs":30}`),
		},
		{
			Title:       "番茄鸡蛋面",
			Description: "酸甜开胃，简单快手，适合晚餐轻食",
			CategoryID:  5,
			CookTime:    15,
			Difficulty:  "简单",
			PeopleCount: 1,
			Taste:       "酸甜",
			HealthTags:  model.JSON(`["快手","清淡"]`),
			Ingredients: model.JSON(`[{"name":"面条","amount":"150g"},{"name":"番茄","amount":"1个"},{"name":"鸡蛋","amount":"2个"}]`),
			Seasonings:  model.JSON(`[{"name":"盐","amount":"适量"},{"name":"葱","amount":"适量"}]`),
			Steps:       model.JSON(`[{"step":1,"desc":"番茄切块，鸡蛋打散"},{"step":2,"desc":"炒鸡蛋盛出"},{"step":3,"desc":"炒番茄出汁，加水煮开"},{"step":4,"desc":"下面条煮熟，加鸡蛋调味"}]`),
			Nutrition:   model.JSON(`{"calories":420,"protein":18,"fat":10,"carbs":65}`),
		},
		{
			Title:       "鸡蛋三明治",
			Description: "营养早餐，简单快捷",
			CategoryID:  3,
			CookTime:    10,
			Difficulty:  "简单",
			PeopleCount: 1,
			Taste:       "清淡",
			HealthTags:  model.JSON(`["快手","高蛋白"]`),
			Ingredients: model.JSON(`[{"name":"全麦面包","amount":"2片"},{"name":"鸡蛋","amount":"2个"},{"name":"生菜","amount":"2片"},{"name":"番茄","amount":"2片"}]`),
			Seasonings:  model.JSON(`[{"name":"盐","amount":"适量"},{"name":"黑胡椒","amount":"少许"}]`),
			Steps:       model.JSON(`[{"step":1,"desc":"煎蛋"},{"step":2,"desc":"面包铺上生菜、番茄、煎蛋"},{"step":3,"desc":"撒上盐和黑胡椒，盖上另一片面包"}]`),
			Nutrition:   model.JSON(`{"calories":350,"protein":20,"fat":12,"carbs":40}`),
		},
	}

	for _, r := range recipes {
		db.Where("title = ?", r.Title).FirstOrCreate(&r)
	}

	// Seed ingredients
	ingredients := []model.Ingredient{
		{Name: "番茄", Category: "蔬菜"},
		{Name: "黄瓜", Category: "蔬菜"},
		{Name: "土豆", Category: "蔬菜"},
		{Name: "茄子", Category: "蔬菜"},
		{Name: "青椒", Category: "蔬菜"},
		{Name: "白菜", Category: "蔬菜"},
		{Name: "西兰花", Category: "蔬菜"},
		{Name: "胡萝卜", Category: "蔬菜"},
		{Name: "猪肉", Category: "肉类"},
		{Name: "牛肉", Category: "肉类"},
		{Name: "鸡肉", Category: "肉类"},
		{Name: "鸡胸肉", Category: "肉类"},
		{Name: "排骨", Category: "肉类"},
		{Name: "虾", Category: "海鲜"},
		{Name: "鱼", Category: "海鲜"},
		{Name: "鸡蛋", Category: "蛋奶"},
		{Name: "牛奶", Category: "蛋奶"},
		{Name: "盐", Category: "调料"},
		{Name: "生抽", Category: "调料"},
		{Name: "老抽", Category: "调料"},
		{Name: "料酒", Category: "调料"},
		{Name: "醋", Category: "调料"},
		{Name: "糖", Category: "调料"},
		{Name: "冰糖", Category: "调料"},
	}

	for _, ing := range ingredients {
		db.Where("name = ?", ing.Name).FirstOrCreate(&ing)
	}

	// Seed admin user
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := model.AdminUser{
		Username:     "admin",
		PasswordHash: string(hash),
		Role:         "admin",
	}
	db.Where("username = ?", "admin").FirstOrCreate(&admin)

	// Seed banners
	banners := []model.Banner{
		{Title: "今日推荐", Image: "/uploads/banner1.jpg", LinkType: "recipe", LinkValue: "1", Sort: 1, Status: 1},
		{Title: "一周菜单", Image: "/uploads/banner2.jpg", LinkType: "url", LinkValue: "/week-menu", Sort: 2, Status: 1},
	}
	for _, b := range banners {
		db.Where("title = ?", b.Title).FirstOrCreate(&b)
	}

	fmt.Println("Seed data inserted successfully!")
}
