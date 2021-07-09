export class Reconnector {
    constructor(config) {
        config = config || {}
        this.timeout = config.initialTimeout || 1000
        this.maxRetries = config.retries || 10
        this.retriesOut = 0
        this.fn = config.fn || (t => t * 2)
    }

    reset() {
        this.timeout = 100
        this.retriesOut = 0
    }

    next() {
        return new Promise((ok, err) => {
            if (this.retriesOut >= this.maxRetries) {
                err('retries are out')
                return
            }
            this.retriesOut += 1
            setTimeout(ok, this.timeout)
            this.timeout = this.fn(this.timeout)
        })
    }
}