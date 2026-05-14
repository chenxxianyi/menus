Component({
  properties: {
    recipe: { type: Object, value: {} },
    showFavorite: { type: Boolean, value: true },
    compact: { type: Boolean, value: false }
  },

  methods: {
    onTap() {
      this.triggerEvent('tap', { recipe: this.data.recipe })
    },

    onFavorite() {
      this.triggerEvent('favorite', { recipe: this.data.recipe })
    }
  }
})
