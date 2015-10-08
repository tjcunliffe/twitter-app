import React from 'react'
import { ProgressBar, Input, ButtonInput } from 'react-bootstrap'
import Highlight from 'react-highlight'
var pd = require('pretty-data').pd;


const TweetsResultsComponent = React.createClass({
    displayName: "TweetsResultsComponent",

    getInitialState() {
        return {
            url: this.props.url,
            data: "",
            loading: true
        }
    },

    componentDidMount() {

        let that = this;

        $.ajax({
            type: "GET",
            dataType: "json",
            url: this.state.url,
            success: function (data) {
                console.log(data);

                if (that.isMounted()) {
                    that.setState({
                        data: data,
                        loading: false
                    });
                }
            }
        });
    },

    render: function () {
        // loading bar
        const progressInstance = (
            <ProgressBar active now={87}/>
        );


        if (this.state.loading) {
            return progressInstance
        } else {
            // since we are not loading, proceeding to data
            var value = this.state.data;
            let prettyfied = pd.json(value);
            let preInstance = (
                <Highlight>
                    {prettyfied}
                </Highlight>
            );
            return <div> {preInstance} </div>
        }


    }
});


const ButtonInputQueryComponent = React.createClass({
    getInitialState() {
        return {
            disabled: true,
            style: null,
            ss: true
        };
    },

    validationState() {
        let length = this.refs.input.getValue().length;
        let scenarioSession = this.refs.scenarioSession.getValue().length;
        let backend = this.refs.backend.getValue();
        let style = 'danger';

        this.state.ss = backend == "https://api.twitter.com";

        if(backend != "http://localhost:8300"){
            if (length > 0) style = 'success';
        } else {
            if (length >0 && scenarioSession > 0) style = 'success';
        }

        //else if (length > 5) style = 'warning';

        let disabled = style !== 'success';

        return {style, disabled};
    },

    handleChange() {
        this.setState(this.validationState());
    },

    handleSubmit(e) {

        e.preventDefault();
        let query = this.refs.input.getValue();
        let backend = this.refs.backend.getValue();
        let scenarioSession = this.refs.scenarioSession.getValue();


        // creating final href
        var href = "/query?q=" + query + "&backend=" + backend + "&ss=" + scenarioSession;
        //var href = "/rand?q=" + query;
        console.log("preparing query");
        console.log(href);

        // unmounting current results
        React.unmountComponentAtNode(document.getElementById('results'));
        //// mounting results
        React.render(<TweetsResultsComponent url={href}/>, document.getElementById("results"))
    },

    render() {
        return (
            <div>
                <form onSubmit={this.handleSubmit}>

                    <Input type="select" label="External API URI" placeholder="select" ref="backend" onChange={this.handleChange}>
                        <option value="https://api.twitter.com">Twitter</option>
                        <option value="http://localhost:8300">Mirage Proxy</option>
                    </Input>

                    <Input type="text" ref="scenarioSession" placeholder="Provide 'scenario:session' here"
                           disabled={this.state.ss}
                           onChange={this.handleChange}/>

                    <Input type="text" ref="input" placeholder="Enter your query here.." onChange={this.handleChange}/>

                    <ButtonInput type="submit" value="Submit"
                                 bsStyle={this.state.style} bsSize="small"
                                 disabled={this.state.disabled}/>
                </form>
            </div>
        );
    }
});

React.render(<ButtonInputQueryComponent />, document.getElementById("app"));