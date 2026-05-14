Page({
  data: {
    meals: [
      { key: 'breakfast', name: '早餐', emoji: '🌅', count: 3 },
      { key: 'lunch', name: '午餐', emoji: '☀', count: 5 },
      { key: 'dinner', name: '晚餐', emoji: '🌙', count: 4 }
    ],
    categories: [
      { key: 'home', name: '家常菜', emoji: '🏠', bg: 'var(--mint-bg)' },
      { key: 'quick', name: '快手菜', emoji: '⚡', bg: 'var(--lemon-bg)' },
      { key: 'diet', name: '减脂餐', emoji: '🥗', bg: 'var(--lavender-bg)' },
      { key: 'kids', name: '儿童餐', emoji: '👶', bg: 'var(--rose-bg)' },
      { key: 'elder', name: '老人餐', emoji: '👴', bg: 'var(--sky-bg)' },
      { key: 'breakfast', name: '早餐', emoji: '🌅', bg: 'var(--peach-bg)' },
      { key: 'soup', name: '汤粥', emoji: '🍲', bg: 'var(--mint-bg)' },
      { key: 'baking', name: '烘焙', emoji: '🧁', bg: 'var(--lemon-bg)' }
    ],
    hotRecipes: [
      { id: 1, name: '番茄炒蛋', time: 15, difficulty: '简单', tags: ['家常菜', '快手'], cover: '', isFavorite: false },
      { id: 2, name: '蒜蓉西兰花', time: 12, difficulty: '简单', tags: ['家常菜', '减脂'], cover: '', isFavorite: true },
      { id: 3, name: '宫保鸡丁', time: 25, difficulty: '中等', tags: ['家常菜', '辣'], cover: '', isFavorite: false },
      { id: 4, name: '减脂鸡胸沙拉', time: 10, difficulty: '简单', tags: ['减脂餐'], cover: '', isFavorite: true }
    ]
  },

  goSearch() {
    wx.navigateTo({ url: '/pages/recipes/list?focus=true' })
  },

  goRecommend() {
    wx.navigateTo({ url: '/pages/recommend/index' })
  },

  goMealRecommend(e) {
    const meal = e.currentTarget.dataset.meal
    wx.navigateTo({ url: `/pages/recommend/index?meal=${meal}` })
  },

  goCategory(e) {
    const key = e.currentTarget.dataset.key
    wx.navigateTo({ url: `/pages/recipes/list?category=${key}` })
  },

  goDetail(e) {
    const id = e.detail.recipe.id
    wx.navigateTo({ url: `/pages/recipes/detail?id=${id}` })
  },

  toggleFavorite(e) {
    const id = e.detail.recipe.id
    const recipes = this.data.hotRecipes
    const recipe = recipes.find(r => r.id === id)
    if (recipe) {
      recipe.isFavorite = !recipe.isFavorite
      this.setData({ hotRecipes: recipes })
      wx.showToast({ title: recipe.isFavorite ? '已收藏' : '已取消', icon: 'none' })
    }
  }
})
