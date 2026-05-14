Page({
  data: {
    tab: 'recipes',
    favRecipes: [
      { id: 1, name: '蒜蓉西兰花', time: 12, difficulty: '简单', tags: ['家常菜'], cover: '', isFavorite: true },
      { id: 2, name: '减脂鸡胸沙拉', time: 10, difficulty: '简单', tags: ['减脂餐'], cover: '', isFavorite: true }
    ],
    favMenus: [
      { id: 1, name: '营养午餐组合', dishes: 3, createdAt: '2025-04-12' }
    ],
    history: [
      { id: 3, name: '宫保鸡丁', time: 25, difficulty: '中等', tags: ['家常菜'], cover: '', isFavorite: false },
      { id: 4, name: '紫薯燕麦粥', time: 20, difficulty: '简单', tags: ['早餐'], cover: '', isFavorite: false }
    ]
  },

  setTab(e) { this.setData({ tab: e.currentTarget.dataset.val }) },

  goDetail(e) {
    wx.navigateTo({ url: `/pages/recipes/detail?id=${e.detail.recipe.id}` })
  },

  removeFav(e) {
    const id = e.detail.recipe.id
    this.setData({ favRecipes: this.data.favRecipes.filter(r => r.id !== id) })
    wx.showToast({ title: '已取消收藏', icon: 'none' })
  },

  goExplore() {
    wx.switchTab({ url: '/pages/recipes/list' })
  }
})
