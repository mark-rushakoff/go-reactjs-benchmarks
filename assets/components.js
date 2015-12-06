var ChildComponent = React.createClass({
  render: function() {
    return React.createElement("span", {className: this.props.className}, this.props.text);
  }
});

var ParentComponent = React.createClass({
  render: function() {
    return React.createElement("div", {className: "the-parent"},
      React.createElement(ChildComponent, {className: "child-1", text: "First child"}),
      React.createElement(ChildComponent, {className: "child-2", text: "Second child"}),
      React.createElement(ChildComponent, {className: "child-3", text: "Third child"})
    );
  }
});
