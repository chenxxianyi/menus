Page({
  data: {
    groups: [
      {
        category: '蔬菜',
        items: [
          { name: '番茄', amount: '2个', checked: false },
          { name: '西兰花', amount: '1颗', checked: true },
          { name: '葱', amount: '适量', checked: false }
        ]
      },
      {
        category: '肉类',
        items: [
          { name: '鸡蛋', amount: '6个', checked: false }
        ]
      },
      {
        category: '调料',
        items: [
          { name: '盐', amount: '适量', checked: false },
          { name: '食用油', amount: '适量', checked: false }
        ]
      }
    ]
  },

  toggleCheck(e) {
    const { cat, name } = e.currentTarget.dataset
    const groups = this.data.groups
    const group = groups.find(g => g.category === cat)
    if (group) {
      const item = group.items.find(i => i.name === name)
      if (item) {
        item.checked = !item.checked
        this.setData({ groups })
      }
    }
  },

  removeItem(e) {
    const { cat, name } = e.currentTarget.dataset
    const groups = this.data.groups
    const group = groups.find(g => g.category === cat)
    if (group) {
      group.items = group.items.filter(i => i.name !== name)
      if (!group.items.length) {
        this.setData({ groups: groups.filter(g => g.category !== cat) })
      } else {
        this.setData({ groups })
      }
    }
  },

  addItem(e) {
    const val = e.detail.value.trim()
    if (!val) return
    const groups = this.data.groups
    let group = groups.find(g => g.category === '其他')
    if (!group) {
      group = { category: '其他', items: [] }
      groups.push(group)
    }
    group.items.push({ name: val, amount: '1份', checked: false })
    this.setData({ groups })
  },

  clearAll() {
    wx.showModal({
      title: '确认清空',
      content: '确定要清空购物清单吗？',
      success: (res) => {
        if (res.confirm) this.setData({ groups: [] })
      }
    })
  },

  goRecommend() {
    wx.navigateTo({ url: '/pages/recommend/index' })
  }
})
