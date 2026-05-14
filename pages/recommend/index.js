const { meals, tastes, goals } = require('../../utils/constants')

Page({
  data: {
    meals,
    tastes,
    goals,
    meal: '',
    people: 2,
    taste: '',
    goal: '',
    haveIngredients: [],
    avoidIngredients: [],
    generating: false
  },

  onLoad(options) {
    if (options.meal) this.setData({ meal: options.meal })
  },

  setMeal(e) { this.setData({ meal: e.currentTarget.dataset.val }) },
  setTaste(e) { this.setData({ taste: e.currentTarget.dataset.val }) },
  setGoal(e) { this.setData({ goal: e.currentTarget.dataset.val }) },

  changePeople(e) {
    const delta = parseInt(e.currentTarget.dataset.delta)
    const people = Math.max(1, Math.min(20, this.data.people + delta))
    this.setData({ people })
  },

  addHave(e) {
    const val = e.detail.value.trim()
    if (!val || this.data.haveIngredients.includes(val)) return
    this.setData({ haveIngredients: [...this.data.haveIngredients, val] })
  },

  removeHave(e) {
    const val = e.currentTarget.dataset.val
    this.setData({ haveIngredients: this.data.haveIngredients.filter(i => i !== val) })
  },

  addAvoid(e) {
    const val = e.detail.value.trim()
    if (!val || this.data.avoidIngredients.includes(val)) return
    this.setData({ avoidIngredients: [...this.data.avoidIngredients, val] })
  },

  removeAvoid(e) {
    const val = e.currentTarget.dataset.val
    this.setData({ avoidIngredients: this.data.avoidIngredients.filter(i => i !== val) })
  },

  generate() {
    if (!this.data.meal) {
      wx.showToast({ title: '请选择用餐时间', icon: 'none' })
      return
    }
    this.setData({ generating: true })
    setTimeout(() => {
      this.setData({ generating: false })
      wx.navigateTo({ url: '/pages/recommend/result' })
    }, 1500)
  }
})
