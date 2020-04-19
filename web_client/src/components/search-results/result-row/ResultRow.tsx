import * as React from 'react'
import {Component, ReactNode} from 'react'
import {SearchMatch} from "../../../services/index/SearchResult";
import {Languages} from "./languages";
import './ResultRow.scss'
import * as hljs from "highlight.js";

interface Props {
    match: SearchMatch
}

export class ResultRow extends Component<Props, {}> {
    private el: React.RefObject<HTMLDivElement>;

    constructor(props: Props) {
        super(props);
        this.el = React.createRef();
    }

    public componentDidMount() {
        this.highlightCode();
    }

    public componentDidUpdate() {
        this.highlightCode();
    }

    public highlightCode() {
        if (!this.el.current) {
            console.error("Cannot highlight, ref is null");
            return;
        }

        const nodes = this.el.current.querySelectorAll('pre code');
        for (let i = 0; i < nodes.length; i++) {
            hljs.highlightBlock(nodes[i])
        }
    }

    public render(): ReactNode {
        const match = this.props.match;
        const language = this.getLanguageFromFilePath(match.File.Path);
        return (
            <div className="card mb-2" ref={this.el}>
                <div className="card-header gs-metadata">
                    <div>Path: {match.File.Path}</div>
                    <div className="flex-grow-1">&nbsp;</div>
                    <div>Commit: {match.File.Commit.substr(0, 15)}</div>
                    {/*TODO: add date and more*/}
                </div>
                <div className="card-body">
                    <pre>
                        <code className={language}
                              dangerouslySetInnerHTML={{__html: match.Fragments.join("\n[...]\n")}}/>
                    </pre>
                </div>
            </div>
        );
    }

    private getLanguageFromFilePath(path: string): string {
        const extension = path.split('.').pop();
        if (!extension) {
            return '';
        }

        const language = Languages.find(lang => lang.extensions.indexOf(extension) !== -1);
        return language?.name || '';
    }
}
