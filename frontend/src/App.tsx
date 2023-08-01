import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Search} from "../wailsjs/go/main/App";

function App() {
    const [name, setName] = useState('');
    const updateName = (e: any) => setName(e.target.value);

    function greet() {
        Search(name);
    }

    return (
        <div id="App">
            <div id="input" className="search-box">
                <input id="name" className="search-input" onChange={updateName} name="input" type="text" spellCheck="false"/>
            </div>

            <div id="results" className="results-div">
                <h1>No Results found</h1>
            </div>
        </div>
    )
}

export default App
