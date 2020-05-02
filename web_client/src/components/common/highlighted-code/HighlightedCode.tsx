import * as React from 'react'
import {Component, ReactNode} from 'react'
import {languageFromFileName} from "../../../core/utils/languages";
import * as hljs from "highlight.js";
import {Logger} from "../../../core/utils/logger";
import './HighlightedCode.scss'

interface Props {
    content: string
    fileName: string
}

export class HighlightedCode extends Component<Props, {}> {
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

    public render(): ReactNode {
        const content = this.props.content;
        const language = languageFromFileName(this.props.fileName);
        return (
            <div ref={this.el}>
                <pre>
                    <code className={language} dangerouslySetInnerHTML={{__html: this.props.content}}/>
                </pre>
            </div>
        );
    }

    public highlightCode() {
        if (!this.el.current) {
            Logger.error("Cannot highlight code, ref is null");
            return;
        }

        const nodes = this.el.current.querySelectorAll('pre code');
        for (let i = 0; i < nodes.length; i++) {
            hljs.highlightBlock(nodes[i])
        }
    }

}
