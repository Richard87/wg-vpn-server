import {useState} from "react"
import {useQuery} from "react-query"

function App() {
    const [name, setName] = useState("")
    const [ip, setIp] = useState("")
    const [publicKey, setPublicKey] = useState("")
    const {data, isSuccess, refetch} = useQuery(["clients"], {
        queryFn: () => fetch(`${process.env.REACT_APP_API_SERVER}/clients`).then(res => res.json())
    })


    const onSubmit = (events) => {
        events.preventDefault()

        fetch(`${process.env.REACT_APP_API_SERVER}/clients`, {
            method: "POST",
            body: JSON.stringify({name, ip, publicKey})
        }).then(() => refetch())

        return false
    }

    return (
        <div className="App">
            <ul>
                {isSuccess && data.map(client => <li key={client.id}>{client.name}</li>)}
            </ul>
            <form onSubmit={onSubmit}>
                <input value={name} onChange={e => setName(e.target.value)} name="name"/>

                <input value={ip} onChange={e => setIp(e.target.value)} name="ip"/>
                <input value={publicKey} onChange={e => setPublicKey(e.target.value)} name="publicKey"/>
                <button type={"submit"}>Lagre klient</button>
            </form>
        </div>
    );
}

export default App;
