package worker

// Middleware wraps a JobHandler (like HTTP middleware)
type Middleware func(JobHandler) JobHandler

// Chain builds a pipeline of middlewares around a JobHandler
func Chain(h JobHandler, middlewares ...Middleware) JobHandler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        if middlewares[i] == nil {
            panic("nil middleware in Chain")
        }
        h = middlewares[i](h)
    }
    return h
}
