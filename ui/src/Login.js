import {useState} from "react";
import {useNavigate} from "react-router-dom";
import {MDBBtn, MDBCard, MDBCardBody, MDBCardHeader, MDBIcon, MDBInput} from "mdbreact";
import {authFetch} from "./index";

export default function Login () {
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState(null)
    const nav = useNavigate()

    const onLogin = () => {
        authFetch(`authenticate`, {
            body: JSON.stringify({username,password}),
            method: "POST"
        })
            .then(({token}) => window.localStorage.setItem("jwt", token))
            .then(() => nav("/dashboard"))
            .catch(err => setError("Unknown username and/or password"))
    }

    return <div style={{paddingTop:"2em"}}><MDBCard>
        <MDBCardHeader className="form-header deep-blue-gradient">
            <h3 className="my-3">
                <MDBIcon icon="lock"/> Login:
            </h3>
        </MDBCardHeader>
        <MDBCardBody>

            <MDBInput
                label="Type your username"
                onChange={e => setUsername(e.target.value)}
                value={username}
                icon="user"
                type="text"
                required
            />
            <MDBInput
                label="Type your password"
                onChange={e => setPassword(e.target.value)}
                value={password}
                icon="lock"
                type="password"
                required
            />
            {error && (
                <div className="text-danger text-center mt-4">
                    {error}
                </div>
            )}


            <div className="text-center mt-4">
                <MDBBtn color="light-blue" className="mb-3" type="button" onClick={onLogin}>
                    Login
                </MDBBtn>
            </div>
        </MDBCardBody>
    </MDBCard></div>
}
