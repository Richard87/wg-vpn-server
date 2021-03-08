import {useQuery} from "react-query"

function App() {

  const {data, isSuccess} = useQuery(["clients"], {queryFn: () => fetch( `${process.env.REACT_APP_API_SERVER}/clients`).then(res => res.json())})

  return (
    <div className="App">
        <ul>
          {isSuccess && data.map(client => <li key={client.id}>{client.name}</li>)}
        </ul>
    </div>
  );
}

export default App;
