import * as React from 'react'
import {Component, ReactNode} from 'react'
import {NothingFound} from "./nothing-found/NothingFound";
import {ResultRow} from "./result-row/ResultRow";
import {SearchResult} from "../../services/index/SearchResult";

import './SearchResults.scss'

export interface Props {
    searchResponse: SearchResult | undefined
}

export class SearchResults extends Component<Props, {}> {
    public render(): ReactNode {
        const matches = this.props.searchResponse?.Matches || [];
        const results = matches.map((match, idx) => <ResultRow match={match} key={idx}/>);
        const nothingFound = matches.length < 1 && !!this.props.searchResponse;
        return (
            <div className="gs-search-results">
                {nothingFound && <NothingFound/>}
                {matches.length > 0 && results}
            </div>
        );
    }
}
