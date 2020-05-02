import * as React from 'react'
import {Component, ReactNode} from 'react'
import {servicesMap} from "../../core/service-factory";
import {RouteComponentProps} from 'react-router-dom';
import {IndexedFile} from "../../core/index/SearchResult";
import {HighlightedCode} from "../common/highlighted-code/HighlightedCode";

type FileParams = {
    fileId: string
}

type Props = RouteComponentProps<FileParams>

interface State {
    file?: IndexedFile
}

export class FileView extends Component<Props, State> {

    private services = servicesMap();

    constructor(props: Props) {
        super(props);
        this.state = {};
    }

    public componentDidMount(): void {
        this.services.repository.getDocumentById(this.props.match.params.fileId)
            .then(res => this.setState({file: res}));
    }

    public render(): ReactNode {
        const file = this.state.file;
        if (!file) {
            return <div/>
        }
        return (
            <>
                <div className="card mb-2">
                    <div className="card-header">
                        <div>File: {file.Name}</div>
                        <div>Path: {file.Path}</div>
                        <div className={"mr-4"}>First seen: {file.Date}</div>
                        <div>Commit: {file.Commit.substr(0, 15)}</div>
                    </div>
                    <div className="card-body">
                        <HighlightedCode content={file.Content} fileName={file.Name}/>
                    </div>
                </div>
                <div className="flex-grow-1">&nbsp;</div>
            </>
        );
    }

}
