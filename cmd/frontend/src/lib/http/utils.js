// Loads defaults and merges any custom options the caller provides.
function inheritOptions(from, opts) {
    return opts ? opts = { ...from, ...opts } : opts
}

function filterObjectFields(input, fields) {
    let input_stripped = { ...input }
    let data = {}
    for (let key in input_stripped) {
        if (fields.includes(key)) {
            data[key] = input_stripped[key]
        }
    }
    return data
}

// Wrap the error responses to be more user friendly, but still provide the option
// of inspection of the response.
function formatError(error) {
    if (error.response) {
        if (error.response.data.errorResponse && error.response.data.errorResponse.errorType) {
            error.message = error.response.data.errorResponse.errorType
        } else if (error.response.data.Error) {
            error.message = error.response.data.Error
        } else if (error.response.data.error) {
            error.message = error.response.data.error
        } else if (error.response.data.message) {
            error.message = error.response.data.message
        } else {
            error.message = error.response.statusText
        }
    }
    return error
}

export default {
    inheritOptions,
    filterObjectFields,
    formatError
}
