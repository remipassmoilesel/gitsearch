import * as React from 'react'
import {Component, ReactNode} from 'react'
import {SearchBar, SearchEvent} from "./search-bar/SearchBar";
import {ControlBar} from "./control-bar/ControlBar";
import {SearchResults} from "./search-results/SearchResults";
import {SearchResult} from "../../core/index/SearchResult";
import {servicesMap} from "../../core/service-factory";
import {RouteComponentProps} from 'react-router-dom';

interface State {
    lastResponse: SearchResult | undefined;
}

interface Params {
    query?: string
}

type Props = RouteComponentProps<Params>

export class SearchView extends Component<Props, State> {

    private services = servicesMap();

    constructor(props: Props) {
        super(props);
        this.state = {
            lastResponse: undefined
        }
    }

    public componentDidMount(): void {
        const query = this.props.match.params.query;
        if (query) {
            this.onSearch({query: decodeURI(query)})
        }
    }

    public render(): ReactNode {
        return (
            <>
                <SearchBar onSearch={this.onSearch}/>
                <ControlBar/>
                <SearchResults searchResponse={this.state.lastResponse}/>
            </>
        );
    }

    private onSearch = (ev: SearchEvent) => {
        this.props.history.push(`/search/${encodeURI(ev.query)}`)
        this.services.index.search(ev.query)
            .then(res => {
                this.setState({lastResponse: res});
                this.services.uiMessages.info(`Search took ${res.TookMs} ms`);
            })
            .catch(err => this.services.error.onError(err))
    }
}
