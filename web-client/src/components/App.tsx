import * as React from 'react'
import {Component, ReactNode} from 'react'
import {Banner} from "./banner/Banner";
import {SearchResults} from "./search-results/SearchResults";
import {SearchBar, SearchEvent} from "./search-bar/SearchBar";
import {servicesMap} from "../services/service-factory";
import {SearchResponse} from "../services/SearchResponse";

import './App.scss'

interface State {
    lastResponse: SearchResponse | undefined;
}

export class App extends Component<{}, State> {

    private services = servicesMap();

    constructor(props: {}) {
        super(props);
        this.state = {
            lastResponse: undefined
        }
    }

    public render(): ReactNode {
        return (
            <div className="gs-container">
                <Banner/>
                <SearchBar onSearch={this.onSearch}/>
                <SearchResults searchResponse={this.state.lastResponse}/>
            </div>
        );
    }

    private onSearch = (ev: SearchEvent) => {
        console.log(`Searching: ${ev}`);
        this.services.search.search(ev.query)
            .then(res => this.setState({lastResponse: res}))
            .catch(err => this.services.error.onError(err))
    }
}
