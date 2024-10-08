!(function (e) {
  function t(n) {
    if (i[n]) return i[n].exports
    var s = (i[n] = { i: n, l: !1, exports: {} })
    return e[n].call(s.exports, s, s.exports, t), (s.l = !0), s.exports
  }
  var i = {}
  return (
    (t.m = e),
    (t.c = i),
    (t.i = function (e) {
      return e
    }),
    (t.d = function (e, t, i) {
      Object.defineProperty(e, t, { configurable: !1, enumerable: !0, get: i })
    }),
    (t.n = function (e) {
      var i =
        e && e.__esModule
          ? function () {
              return e.default
            }
          : function () {
              return e
            }
      return t.d(i, 'a', i), i
    }),
    (t.o = function (e, t) {
      return Object.prototype.hasOwnProperty.call(e, t)
    }),
    (t.p = '/dist/'),
    t((t.s = 3))
  )
})([
  function (e, t) {
    'use strict'
    Object.defineProperty(t, '__esModule', { value: !0 })
    ;(t.$ = OpenSeadragon),
      (t.DIRECTION_HORIZONTAL = Symbol('horizontal')),
      (t.DIRECTION_VERTICAL = Symbol('vertical'))
  },
  function (e, t, i) {
    'use strict'
    function n(e, t, i, n, s) {
      var d = n === a.DIRECTION_HORIZONTAL ? 'horizontal' : 'vertical',
        l = o(),
        u = l.find(function (t) {
          return t.id === e
        })
      u
        ? ((u.x = t), (u.y = i), (u.rotation = s))
        : l.push({ id: e, x: t, y: i, direction: d, rotation: s }),
        r(l)
    }
    function s(e) {
      var t = o()
      r(
        t.filter(function (t) {
          return t.id !== e
        })
      )
    }
    function o() {
      var e = window.sessionStorage.getItem(l)
      return e ? JSON.parse(e) : []
    }
    function r(e) {
      window.sessionStorage.setItem(l, JSON.stringify(e))
    }
    Object.defineProperty(t, '__esModule', { value: !0 })
    var a = i(0),
      d = !1,
      l = 'openseadragon-guides'
    t.default = { addGuide: n, deleteGuide: s, getGuides: o, useStorage: d }
  },
  function (e, t, i) {
    'use strict'
    function n(e) {
      return e && e.__esModule ? e : { default: e }
    }
    function s(e, t) {
      if (!(e instanceof t))
        throw new TypeError('Cannot call a class as a function')
    }
    function o(e, t) {
      var i = document.createElement('div')
      switch (((i.id = 'osd-guide-' + t), i.classList.add('osd-guide'), e)) {
        case d.DIRECTION_HORIZONTAL:
          i.classList.add('osd-guide-horizontal')
          break
        case d.DIRECTION_VERTICAL:
          i.classList.add('osd-guide-vertical')
          break
        default:
          throw new Error('Invalid guide direction')
      }
      return i
    }
    function r() {
      var e = document.createElement('div')
      return e.classList.add('osd-guide-line'), e
    }
    Object.defineProperty(t, '__esModule', { value: !0 }), (t.Guide = void 0)
    var a = (function () {
        function e(e, t) {
          for (var i = 0; i < t.length; i++) {
            var n = t[i]
            ;(n.enumerable = n.enumerable || !1),
              (n.configurable = !0),
              'value' in n && (n.writable = !0),
              Object.defineProperty(e, n.key, n)
          }
        }
        return function (t, i, n) {
          return i && e(t.prototype, i), n && e(t, n), t
        }
      })(),
      d = i(0),
      l = i(1),
      u = n(l)
    t.Guide = (function () {
      function e(t) {
        var i = t.clickHandler,
          n = t.direction,
          a = void 0 === n ? d.DIRECTION_HORIZONTAL : n,
          l = t.id,
          u = void 0 === l ? Date.now() : l,
          h = t.plugin,
          c = t.rotation,
          p = void 0 === c ? 0 : c,
          v = t.viewer,
          g = t.x,
          f = t.y
        s(this, e),
          (this.viewer = v),
          (this.plugin = h),
          (this.direction = a),
          (this.rotation = p),
          (this.id = u),
          (this.point = this.viewer.viewport.getCenter()),
          (this.point.x = g ? g : this.point.x),
          (this.point.y = f ? f : this.point.y),
          (this.elem = o(this.direction, this.id)),
          (this.line = r()),
          this.elem.appendChild(this.line),
          (this.overlay = new d.$.Overlay(this.elem, this.point)),
          this.draw(),
          this.saveInStorage(),
          i && this.plugin.allowRotation && (this.clickHandler = i),
          this.plugin.allowRotation && this.rotate(this.rotation),
          this.addHandlers()
      }
      return (
        a(e, [
          {
            key: 'addHandlers',
            value: function () {
              ;(this.tracker = new d.$.MouseTracker({
                element: this.elem,
                clickTimeThreshold: this.viewer.clickTimeThreshold,
                clickDistThreshold: this.viewer.clickDistThreshold,
                dragHandler: this.dragHandler.bind(this),
                dragEndHandler: this.dragEndHandler.bind(this),
                dblClickHandler: this.remove.bind(this),
              })),
                this.clickHandler &&
                  (this.tracker.clickHandler = this.clickHandler.bind(this)),
                this.viewer.addHandler('open', this.draw.bind(this)),
                this.viewer.addHandler('animation', this.draw.bind(this)),
                this.viewer.addHandler('resize', this.draw.bind(this)),
                this.viewer.addHandler('rotate', this.draw.bind(this))
            },
          },
          {
            key: 'dragHandler',
            value: function (e) {
              var t = this.viewer.viewport.deltaPointsFromPixels(e.delta, !0)
              switch (this.direction) {
                case d.DIRECTION_HORIZONTAL:
                  this.point.y += t.y
                  break
                case d.DIRECTION_VERTICAL:
                  this.point.x += t.x
              }
              this.elem.classList.add('osd-guide-drag'), this.draw()
            },
          },
          {
            key: 'dragEndHandler',
            value: function () {
              this.elem.classList.remove('osd-guide-drag'), this.saveInStorage()
            },
          },
          {
            key: 'draw',
            value: function () {
              return (
                this.point &&
                  (this.overlay.update(this.point),
                  this.overlay.drawHTML(
                    this.viewer.drawer.container,
                    this.viewer.viewport
                  )),
                this
              )
            },
          },
          {
            key: 'remove',
            value: function () {
              return (
                this.viewer.removeHandler('open', this.draw.bind(this)),
                this.viewer.removeHandler('animation', this.draw.bind(this)),
                this.viewer.removeHandler('resize', this.draw.bind(this)),
                this.viewer.removeHandler('rotate', this.draw.bind(this)),
                this.overlay.destroy(),
                (this.point = null),
                u.default.deleteGuide(this.id),
                this.plugin.allowRotation && this.plugin.closePopup(),
                this
              )
            },
          },
          {
            key: 'rotate',
            value: function (e) {
              if (parseFloat(e)) {
                switch (this.direction) {
                  case d.DIRECTION_HORIZONTAL:
                    ;(this.line.style.webkitTransform =
                      'rotateZ(' + e + 'deg)'),
                      (this.line.style.transform = 'rotateZ(' + e + 'deg)')
                    break
                  case d.DIRECTION_VERTICAL:
                    ;(this.line.style.webkitTransform =
                      'rotateZ(' + e + 'deg)'),
                      (this.line.style.transform = 'rotateZ(' + e + 'deg)')
                }
                this.rotation = e
              } else
                (this.line.style.webkitTransform = ''),
                  (this.line.style.transform = ''),
                  (this.rotation = 0)
              this.saveInStorage()
            },
          },
          {
            key: 'saveInStorage',
            value: function () {
              u.default.useStorage &&
                u.default.addGuide(
                  this.id,
                  this.point.x,
                  this.point.y,
                  this.direction,
                  this.rotation
                )
            },
          },
        ]),
        e
      )
    })()
  },
  function (e, t, i) {
    'use strict'
    function n(e) {
      return e && e.__esModule ? e : { default: e }
    }
    var s = i(2),
      o = i(0),
      r = i(1),
      a = n(r)
    if (!o.$.version || o.$.version.major < 2)
      throw new Error(
        'This version of OpenSeadragon Guides requires OpenSeadragon version 2.0.0+'
      )
    ;(o.$.Viewer.prototype.guides = function (e) {
      return (
        (this.guidesInstance && !e) ||
          ((e = e || {}),
          (e.viewer = this),
          (this.guidesInstance = new o.$.Guides(e))),
        this.guidesInstance
      )
    }),
      (o.$.Guides = function (e) {
        var t = this
        o.$.extend(
          !0,
          this,
          {
            viewer: null,
            guides: [],
            allowRotation: !1,
            horizontalGuideButton: null,
            verticalGuideButton: null,
            prefixUrl: null,
            removeOnClose: !1,
            useSessionStorage: !1,
            navImages: {
              guideHorizontal: {
                REST: 'guidehorizontal_rest.png',
                GROUP: 'guidehorizontal_grouphover.png',
                HOVER: 'guidehorizontal_hover.png',
                DOWN: 'guidehorizontal_pressed.png',
              },
              guideVertical: {
                REST: 'guidevertical_rest.png',
                GROUP: 'guidevertical_grouphover.png',
                HOVER: 'guidevertical_hover.png',
                DOWN: 'guidevertical_pressed.png',
              },
            },
          },
          e
        ),
          o.$.extend(!0, this.navImages, this.viewer.navImages)
        var i = this.prefixUrl || this.viewer.prefixUrl || '',
          n = this.viewer.buttons && this.viewer.buttons.buttons,
          r = n ? this.viewer.buttons.buttons[0] : null,
          d = r ? r.onFocus : null,
          l = r ? r.onBlur : null
        if (
          ((this.horizontalGuideButton = new o.$.Button({
            element: this.horizontalGuideButton
              ? o.$.getElement(this.horizontalGuideButton)
              : null,
            clickTimeThreshold: this.viewer.clickTimeThreshold,
            clickDistThreshold: this.viewer.clickDistThreshold,
            tooltip:
              o.$.getString('Tooltips.HorizontalGuide') || 'Horizontal guide',
            srcRest: i + this.navImages.guideHorizontal.REST,
            srcGroup: i + this.navImages.guideHorizontal.GROUP,
            srcHover: i + this.navImages.guideHorizontal.HOVER,
            srcDown: i + this.navImages.guideHorizontal.DOWN,
            onRelease: this.createHorizontalGuide.bind(this),
            onFocus: d,
            onBlur: l,
          })),
          (this.verticalGuideButton = new o.$.Button({
            element: this.verticalGuideButton
              ? o.$.getElement(this.verticalGuideButton)
              : null,
            clickTimeThreshold: this.viewer.clickTimeThreshold,
            clickDistThreshold: this.viewer.clickDistThreshold,
            tooltip:
              o.$.getString('Tooltips.VerticalGuide') || 'vertical guide',
            srcRest: i + this.navImages.guideVertical.REST,
            srcGroup: i + this.navImages.guideVertical.GROUP,
            srcHover: i + this.navImages.guideVertical.HOVER,
            srcDown: i + this.navImages.guideVertical.DOWN,
            onRelease: this.createVerticalGuide.bind(this),
            onFocus: d,
            onBlur: l,
          })),
          n &&
            (this.viewer.buttons.buttons.push(this.horizontalGuideButton),
            this.viewer.buttons.element.appendChild(
              this.horizontalGuideButton.element
            ),
            this.viewer.buttons.buttons.push(this.verticalGuideButton),
            this.viewer.buttons.element.appendChild(
              this.verticalGuideButton.element
            )),
          (a.default.useStorage = this.useSessionStorage),
          a.default.useStorage)
        ) {
          var u = a.default.getGuides()
          u.forEach(function (e) {
            var i = new s.Guide({
              viewer: t.viewer,
              direction:
                'horizontal' === e.direction
                  ? o.DIRECTION_HORIZONTAL
                  : o.DIRECTION_VERTICAL,
              rotation: e.rotation,
              id: e.id,
              clickHandler: function () {
                return t.showPopup(i)
              },
              plugin: t,
              x: e.x,
              y: e.y,
            })
            t.guides.push(i)
          })
        }
        this.removeOnClose &&
          this.viewer.addHandler('close', function () {
            t.guides.forEach(function (e) {
              return e.remove()
            }),
              (t.guides = [])
          }),
          this.allowRotation &&
            ((this.popup = this.createRotationPopup()),
            this.viewer.addControl(this.popup, {}),
            (this.popup.style.display = 'none'),
            (this.popupInput = this.popup.querySelector('input')))
      }),
      o.$.extend(o.$.Guides.prototype, {
        createHorizontalGuide: function () {
          var e = this,
            t = new s.Guide({
              viewer: this.viewer,
              plugin: this,
              direction: o.DIRECTION_HORIZONTAL,
              clickHandler: function () {
                return e.showPopup(t)
              },
            })
          this.guides.push(t)
        },
        createVerticalGuide: function () {
          var e = this,
            t = new s.Guide({
              viewer: this.viewer,
              plugin: this,
              direction: o.DIRECTION_VERTICAL,
              clickHandler: function () {
                return e.showPopup(t)
              },
            })
          this.guides.push(t)
        },
        showPopup: function (e) {
          ;(this.popup.style.display = 'block'),
            (this.selectedGuide = e),
            (this.popupInput.value = this.selectedGuide.rotation)
        },
        closePopup: function () {
          ;(this.popup.style.display = 'none'),
            (this.popupInput.value = ''),
            (this.selectedGuide = null)
        },
        createRotationPopup: function () {
          var e = this,
            t = document.createElement('div')
          ;(t.id = 'osd-guide-popup'),
            t.classList.add('osd-guide-popup'),
            (t.style.position = 'absolute'),
            (t.style.bottom = '10px'),
            (t.style.left = '10px')
          var i = document.createElement('form')
          i.classList.add('osd-guide-popup-form'),
            (i.style.display = 'block'),
            (i.style.position = 'relative'),
            (i.style.background = '#fff'),
            (i.style.padding = '10px'),
            t.appendChild(i)
          var n = document.createElement('input')
          ;(n.type = 'text'),
            (n.style.display = 'inline-block'),
            (n.style.width = '50px'),
            (n.style.fontSize = '14px'),
            i.appendChild(n)
          var s = document.createElement('button')
          ;(s.type = 'submit'),
            (s.innerHTML = o.$.getString('Tool.rotate') || 'rotate'),
            (s.style.fontSize = '14px'),
            s.classList.add('osd-guide-rotate-button'),
            s.addEventListener('click', function () {
              e.selectedGuide.rotate(n.value), e.closePopup()
            }),
            i.appendChild(s)
          var r = document.createElement('button')
          return (
            (r.innerHTML = '&times;'),
            (r.title = o.$.getString('Tool.close') || 'close'),
            (r.style.fontWeight = 'bold'),
            (r.style.fontSize = '14px'),
            r.classList.add('osd-guide-close'),
            r.addEventListener('click', function () {
              e.closePopup()
            }),
            i.appendChild(r),
            t
          )
        },
      })
  },
])
//# sourceMappingURL=openseadragon-guides.js.map
