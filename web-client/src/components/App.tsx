import * as React from 'react'
import {Component, ReactNode} from 'react'
import {ContextBanner} from "./context-banner/ContextBanner";
import {SearchResults} from "./search-results/SearchResults";
import {SearchBar, SearchEvent} from "./search-bar/SearchBar";
import {servicesMap} from "../services/service-factory";
import {SearchResult} from "../services/index/SearchResult";
import {ControlBar} from "./control-bar/ControlBar";
import {MessageBanner} from "./control-bar/message-banner/MessageBanner";
import './App.scss'

interface State {
    lastResponse: SearchResult | undefined;
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
            <div className="gs-app-container">
                <ContextBanner/>
                <SearchBar onSearch={this.onSearch}/>
                <ControlBar/>
                <SearchResults searchResponse={this.state.lastResponse}/>
            </div>
        );
    }

    private onSearch = (ev: SearchEvent) => {
        this.services.index.search(ev.query)
            .then(res => {
                this.setState({lastResponse: res})
                this.services.uiMessages.info(`Search took ${res.TookSeconds} seconds`)
            })
            .catch(err => this.services.error.onError(err))
    }
}
