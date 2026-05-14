const { categories, sortOptions } = require('../../utils/constants')

Page({
  data: {
    keyword: '',
    focus: false,
    categories,
    sortOptions,
    activeCategory: '',
    sortBy: 'hot',
    recipes: [
      { id: 1, name: '番茄炒蛋', category: 'home', time: 15, difficulty: '简单', tags: ['家常菜'], cover: '', isFavorite: false },
      { id: 2, name: '蒜蓉西兰花', category: 'home', time: 12, difficulty: '简单', tags: ['家常菜'], cover: '', isFavorite: true },
      { id: 3, name: '宫保鸡丁', category: 'home', time: 25, difficulty: '中等', tags: ['家常菜'], cover: '', isFavorite: false },
      { id: 4, name: '减脂鸡胸沙拉', category: 'diet', time: 10, difficulty: '简单', tags: ['减脂餐'], cover: '', isFavorite: true },
      { id: 5, name: '紫薯燕麦粥', category: 'breakfast', time: 20, difficulty: '简单', tags: ['早餐'], cover: '', isFavorite: false },
      { id: 6, name: '红烧排骨', category: 'home', time: 45, difficulty: '中等', tags: ['家常菜'], cover: '', isFavorite: false },
      { id: 7, name: '酸奶水果杯', category: 'quick', time: 5, difficulty: '简单', tags: ['快手菜'], cover: '', isFavorite: false },
      { id: 8, name: '儿童蔬菜饼', category: 'kids', time: 15, difficulty: '简单', tags: ['儿童餐'], cover: '', isFavorite: false }
    ],
    filteredRecipes: []
  },

  onLoad(options) {
    if (options.focus) this.setData({ focus: true })
    if (options.category) {
      this.setData({ activeCategory: options.category })
    }
    this.filterRecipes()
  },

  onInput(e) {
    this.setData({ keyword: e.detail.value })
    this.filterRecipes()
  },

  onSearch() { this.filterRecipes() },

  setCategory(e) {
    this.setData({ activeCategory: e.currentTarget.dataset.key })
    this.filterRecipes()
  },

  clearCategory() {
    this.setData({ activeCategory: '' })
    this.filterRecipes()
  },

  setSort(e) {
    this.setData({ sortBy: e.currentTarget.dataset.key })
    this.filterRecipes()
  },

  filterRecipes() {
    let list = this.data.recipes.filter(r => {
      const matchCat = !this.data.activeCategory || r.category === this.data.activeCategory
      const matchKey = !this.data.keyword || r.name.includes(this.data.keyword)
      return matchCat && matchKey
    })
    this.setData({ filteredRecipes: list })
  },

  goDetail(e) {
    wx.navigateTo({ url: `/pages/recipes/detail?id=${e.detail.recipe.id}` })
  },

  toggleFavorite(e) {
    const id = e.detail.recipe.id
    const recipes = this.data.recipes
    const r = recipes.find(r => r.id === id)
    if (r) {
      r.isFavorite = !r.isFavorite
      this.setData({ recipes })
      this.filterRecipes()
    }
  }
})
