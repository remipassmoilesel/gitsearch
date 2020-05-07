import * as React from 'react'
import {Component, ReactNode} from 'react'
import {NothingFound} from "./nothing-found/NothingFound";
import {ResultCard} from "./result-card/ResultCard";
import {SearchResult} from "../../../core/index/SearchResult";

import './SearchResults.scss'

export interface Props {
    searchResponse: SearchResult | undefined
}

export class SearchResults extends Component<Props, {}> {
    private viewPortRef: React.RefObject<HTMLDivElement>;

    constructor(props: Props) {
        super(props);
        this.viewPortRef = React.createRef()
    }

    public render(): ReactNode {
        const matches = this.props.searchResponse?.Matches || [];
        const results = matches.map((match, idx) => <ResultCard match={match} key={idx}/>);
        const nothingFound = matches.length < 1 && !!this.props.searchResponse;
        return (
            <div className="gs-search-results" ref={this.viewPortRef}>
                {nothingFound && <NothingFound/>}
                {matches.length > 0 && results}
            </div>
        );
    }

    public componentDidUpdate(prevProps: Readonly<Props>, prevState: Readonly<{}>) {
        if(prevProps.searchResponse !== this.props.searchResponse){
            const viewPort = this.viewPortRef.current
            if(viewPort){
               viewPort.scrollTop = 0;
            }
        }
    }
}
