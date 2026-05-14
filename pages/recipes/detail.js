Page({
  data: {
    recipe: {
      id: 1,
      name: '番茄炒蛋',
      cover: '',
      description: '经典家常菜，酸甜可口，老少皆宜',
      time: 15,
      difficulty: '简单',
      servings: 2,
      isFavorite: false,
      ingredients: [
        { name: '番茄', amount: '2个' },
        { name: '鸡蛋', amount: '3个' },
        { name: '葱', amount: '适量' }
      ],
      seasonings: [
        { name: '盐', amount: '适量' },
        { name: '糖', amount: '1勺' },
        { name: '食用油', amount: '适量' }
      ],
      steps: [
        '番茄洗净切块，葱切葱花',
        '鸡蛋打散，加少许盐搅匀',
        '锅热倒油，倒入蛋液炒至凝固盛出',
        '另起锅倒油，放入番茄翻炒出汁',
        '加糖调味，倒入炒好的鸡蛋翻炒均匀',
        '撒葱花出锅'
      ],
      tips: '番茄可先用热水烫一下去皮，口感更好。蛋液加少许水炒出来更嫩。'
    }
  },

  onLoad(options) {
    // TODO: 根据 id 获取菜谱详情
  },

  toggleFavorite() {
    this.setData({ 'recipe.isFavorite': !this.data.recipe.isFavorite })
    wx.showToast({ title: this.data.recipe.isFavorite ? '已收藏' : '已取消', icon: 'none' })
  },

  addToShopping() {
    wx.showToast({ title: '已加入购物清单', icon: 'success' })
  },

  shareRecipe() {
    wx.showShareMenu({ withShareTicket: true })
  }
})
