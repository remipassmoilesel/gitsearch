import * as React from 'react'
import {Component, ReactNode} from 'react'
import {ContextBanner} from "./common/context-banner/ContextBanner";
import {servicesMap} from "../core/service-factory";
import {HashRouter, Route, Switch} from 'react-router-dom';
import {SearchView} from "./search-view/SearchView";
import {FileView} from "./file-view/FileView";
import {MessageBanner} from "./common/message-banner/MessageBanner";
import './App.scss'

export class App extends Component<{}, {}> {

    private services = servicesMap();

    constructor(props: {}) {
        super(props);
        this.state = {
            lastResponse: undefined
        }
    }

    public render(): ReactNode {
        return (
            <HashRouter>
                <div className="gs-app-container">
                    <ContextBanner/>
                    <Switch>
                        <Route path="/" exact component={SearchView}/>
                        <Route path="/search/:query" exact component={SearchView}/>
                        <Route path="/file/:fileId" component={FileView}/>
                    </Switch>
                    <MessageBanner/>
                </div>
            </HashRouter>
        );
    }

}
