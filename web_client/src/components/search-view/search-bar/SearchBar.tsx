import * as React from 'react'
import {ChangeEvent, Component, KeyboardEvent, ReactNode} from 'react'
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faInfo, faSearch} from '@fortawesome/free-solid-svg-icons'

import './SearchBar.scss'

export interface SearchEvent {
    query: string
}

interface Props {
    onSearch: (ev: SearchEvent) => void
}

interface State {
    query: string
}

export class SearchBar extends Component<Props, State> {

    constructor(props: Props) {
        super(props);
        this.state = {
            query: ""
        }
    }

    public render(): ReactNode {
        return (
            <>
                <div className="input-group gs-search-bar">
                    <input type="text"
                           value={this.state.query}
                           onChange={this.onUserInput}
                           className="form-control"
                           onKeyDown={this.onKeyDown}
                           placeholder="Search ..."/>
                    <div className="input-group-append">
                        <button className="btn btn-primary" onClick={this.onSearch}>
                            <FontAwesomeIcon icon={faSearch}/>
                        </button>
                    </div>
                </div>
                <div className="gs-search-help">
                    <FontAwesomeIcon icon={faInfo}/> You can search <i>term</i>, search <i>"exact terms"</i>, include <i>+terms</i>, exclude <i>-terms</i> etc ...
                    See <a href={"https://blevesearch.com/docs/Query-String-Query/"} target="_blank">query strings</a>.
                </div>
            </>
        );
    }

    private onUserInput = (ev: ChangeEvent<HTMLInputElement>) => {
        this.setState({query: ev.target.value})
    };

    private onSearch = () => {
        this.props.onSearch({query: this.state.query})
    };

    private onKeyDown = (ev: KeyboardEvent<HTMLInputElement>) => {
        if (ev.key === 'Enter') {
            this.props.onSearch({query: this.state.query})
        }
    };
}
