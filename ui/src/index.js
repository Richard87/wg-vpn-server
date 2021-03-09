import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import {QueryClientProvider, QueryClient} from "react-query";
import '@fortawesome/fontawesome-free/css/all.min.css';
import 'bootstrap-css-only/css/bootstrap.min.css';
import 'mdbreact/dist/css/mdb.css';
import "./wireguard"

const queryClient = new QueryClient()

ReactDOM.render(
  <React.StrictMode>
      <QueryClientProvider client={queryClient}>
    <App />
      </QueryClientProvider>
  </React.StrictMode>,
  document.getElementById('root')
);
