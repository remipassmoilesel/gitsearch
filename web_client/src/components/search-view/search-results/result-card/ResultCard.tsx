import * as React from 'react'
import {Component, ReactNode} from 'react'
import {SearchMatch} from "../../../../core/index/SearchResult";
import {Link} from 'react-router-dom';
import {HighlightedCode} from "../../../common/highlighted-code/HighlightedCode";
import {faExternalLinkAlt} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import './ResultCard.scss'

interface Props {
    match: SearchMatch
}

export class ResultCard extends Component<Props, {}> {

    public render(): ReactNode {
        const match = this.props.match;
        return (
            <div className="card mb-2">
                <div className="card-header gs-metadata">
                    <div title={match.File.Path}>
                        <Link to={`/file/${match.File.Hash}`} target={"_blank"}>
                            File: {match.File.Name}
                            <FontAwesomeIcon icon={faExternalLinkAlt} className="ml-2"/>
                        </Link>
                    </div>
                    <div className="flex-grow-1">&nbsp;</div>
                    <div className={"mr-4"}>First seen: {match.File.Date}</div>
                    <div>Commit: {match.File.Commit.substr(0, 15)}</div>
                </div>
                <div className="card-body">
                    <HighlightedCode content={match.Fragments.join("\n[...]\n")} fileName={match.File.Name}/>
                </div>
            </div>
        );
    }

}
