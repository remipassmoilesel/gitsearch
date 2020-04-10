import * as React from 'react'
import {Component, ReactNode} from 'react'
import {NothingFound} from "./nothing-found/NothingFound";
import {ResultRow} from "./result-row/ResultRow";
import {SearchResponse} from "../../services/SearchResponse";

import './SearchResults.scss'

export interface Props {
    searchResponse: SearchResponse | undefined
}

export class SearchResults extends Component<Props, {}> {
    public render(): ReactNode {
        const matches = this.props.searchResponse?.matches || [];
        const results = matches.map((match, idx) => <ResultRow key={idx}/>);
        return (
            <div>
                {matches.length < 1 && <NothingFound/>}
                {matches.length > 0 && results}
            </div>
        );
    }
}
