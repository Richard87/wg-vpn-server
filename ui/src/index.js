import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import {QueryClientProvider, QueryClient} from "react-query";
import '@fortawesome/fontawesome-free/css/all.min.css';
import 'bootstrap-css-only/css/bootstrap.min.css';
import 'mdbreact/dist/css/mdb.css';
import "./wireguard"

export const authFetch = async (url, options = {}) => {
    const {headers = {}, ...restOptions} = options
    let jwt = window.localStorage.getItem("jwt");
    headers['Content-type'] = 'application/json'
    headers['Accept'] = 'application/json'
    if (jwt)
        headers.Authorization = "bearer " + jwt

    let response = await fetch(url, {headers, credentials: "include",...restOptions});

    if (response.status === 403) {
        window.location = "/"
    } else if (response.status >= 400) {
        throw new Error(response.status)
    }

    let contentType = response.headers.get("Content-type") || "text/plain"
    if (contentType.indexOf(";") > 0)
        contentType = contentType.substr(0, contentType.indexOf(";"))

    console.log(contentType)
    if (["application/json", "application/ld+json"].includes(contentType))
        return response.json()

    return response.text()
}


const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            queryFn: ({queryKey}) => authFetch(queryKey)
        }
    }
})

ReactDOM.render(
  <React.StrictMode>
      <QueryClientProvider client={queryClient}>
            <App />
      </QueryClientProvider>
  </React.StrictMode>,
  document.getElementById('root')
);
