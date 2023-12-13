import React from 'react';
import logo from './logo.svg';
import './App.css';
import BasicTable from './component/list';
import {
  QueryClient,
  QueryClientProvider,
} from 'react-query';
import { BrowserRouter } from 'react-router-dom';
import FileUpload from './component/upload';
function App() {
  const queryClient = new QueryClient();
  return (
    <BrowserRouter>
    <div className="App">
    <FileUpload/>
     <QueryClientProvider client={queryClient}>

       <BasicTable/>
       </QueryClientProvider>
    </div>
    </BrowserRouter>
  );
}

export default App;
