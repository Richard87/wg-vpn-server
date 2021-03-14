import {
    MDBBtn,
    MDBCard,
    MDBCardBody,
    MDBCardHeader,
    MDBContainer,
    MDBIcon,
    MDBInput, MDBModalFooter,
    MDBNavbar,
    MDBNavbarBrand,
} from "mdbreact";
import {BrowserRouter, Route, Routes, useNavigate} from "react-router-dom";
import Dashboard from "./Dashboard";
import {useState} from "react";


function App() {


    return (
        <div className="App">
            <MDBNavbar color="indigo" dark expand="md">
                <MDBNavbarBrand>
                    <strong className="white-text">WG VPN Server</strong>
                </MDBNavbarBrand>
            </MDBNavbar>
            <MDBContainer>
                <BrowserRouter>
                    <Routes>
                        <Route path={"/dashboard"} element={<Dashboard />} />
                        <Route path={"/"} element={<Login />} />
                    </Routes>
                </BrowserRouter>
            </MDBContainer>
        </div>
    );
}

export default App;

const Login = () => {
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const nav = useNavigate()

    const onLogin = () => {
        console.log(username, password)
        nav("/dashboard")
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
                group
                type="text"
                validate
                error="wrong"
                success="right"
            />
            <MDBInput
                label="Type your password"
                onChange={e => setPassword(e.target.value)}
                value={password}
                icon="lock"
                group
                type="password"
                validate
            />


            <div className="text-center mt-4">
                <MDBBtn color="light-blue" className="mb-3" type="button" onClick={onLogin}>
                    Login
                </MDBBtn>
            </div>
        </MDBCardBody>
    </MDBCard></div>
}

