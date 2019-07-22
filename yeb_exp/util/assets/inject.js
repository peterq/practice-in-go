setTimeout((function (getEventListeners) {
    return function () {
        var config = {
            frameId: '[frameId]',
            contextId: '[contextId]',
            user: '[u]',
            pass: '[pass]'
        }
        config.contextId *= 1

        function notifyChromedp(type, data) {
            console.debug('__notify__' + JSON.stringify({
                contextId: config.contextId,
                type: type,
                data: data,
                frameId: config.frameId,
            }))
            return true
        }

        window.notifyChromedp = notifyChromedp
        function onece(checkFn, cb, int) {
            int = int || 1
            var timer = setInterval(function () {
                var r = checkFn()
                if (r) {
                    clearInterval(timer)
                    cb(r)
                }
            })
        }
        function oneceDom(checkFn, cb) {
            var called = false
            var fn = function () {
                if (called) return
                var dom = checkFn()
                if (dom) {
                    called = true
                    document.removeEventListener('DOMSubtreeModified', fn)
                    cb(dom)
                }
            }
            document.addEventListener('DOMSubtreeModified', fn)
        }

        onece(function () {
            return window.json_ua
        }, function (ua) {
            notifyChromedp('ua.ok', {ua: ua})
            /*setTimeout(function () {
                location.reload()
            }, 5e3)*/
        })
        window.doSubmit = function (m) {
            // document.querySelector('#ant-render-id-pages_outside_components_mobile_mobile > div > input').value = m
            document.querySelector('#ant-render-id-pages_outside_components_mobile_mobile > div > div').click()
            return true
        }
    }
})(getEventListeners));
true