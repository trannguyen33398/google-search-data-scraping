import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import { Pagination } from '@mui/material';
import { makeStyles } from "@material-ui/core/styles";
import { useQuery } from 'react-query';
import { getListResultScraping } from '../api/getListResultScraping';
import { useQueryString } from '../util/utils';
import { useEffect } from 'react';
const socket = new WebSocket("ws://localhost:9090/ws");
export const useStyles = makeStyles({
    tableList: {
     margin: '15% 8%'
    },
  });
  const LIMIT = 10
export default function BasicTable() {
    const classes = useStyles();

useEffect(() => {
    socket.onopen = function() {
      console.log("WebSocket connection established.");
      
      // Send a message to the server
      socket.send("Hello, server!");
    };
    socket.onmessage = function(event) {
      console.log("Received message:", event.data);
    };
    
  

  },[]);
   
   const queryString: { page?: string } = useQueryString()
    const page =queryString.page  ? Number(queryString.page) :1
  
    const dataQuery = useQuery({
      queryKey: ['scraping', page],
      queryFn: () => {
        const controller = new AbortController()
        setTimeout(() => {
          controller.abort()
        }, 5000)
        return getListResultScraping(page, LIMIT, controller.signal)
      },
      keepPreviousData: true,
      retry: 0
    })
   
  return (
    <div className={classes.tableList}>
    <TableContainer component={Paper} >
      <Table sx={{ minWidth: 650 }} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell align="left">Key word</TableCell>
            <TableCell align="left">Total Link  </TableCell>
            <TableCell align="left">Total Search</TableCell>
            <TableCell align="left">Total Advertised</TableCell>
            <TableCell align="left">Html</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
        {dataQuery.data?.data.map((row) => (
            <TableRow
              key={row.keyword}
              sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
              <TableCell component="th" scope="row">
                {row.keyword}
              </TableCell>
              <TableCell align="left">{row.totalAdvertised}</TableCell>
              <TableCell align="left">{row.totalLink}</TableCell>
              <TableCell align="left">{row.totalSearch}</TableCell>
              <TableCell align="left">{row.html}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Pagination count={10} />
    </TableContainer>
    </div>
  );
}