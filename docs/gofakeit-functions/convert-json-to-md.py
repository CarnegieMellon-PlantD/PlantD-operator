import json
from copy import deepcopy

INPUT_FILENAME = "gofakeit-functions.json"
OUTPUT_FILENAME = "gofakeit-functions.md"

HEADERS = ["Name", "Description", "Type", "Parameters"]
PARAM_HEADERS = ["Name", "Description", "Type", "Optional", "Default", "Options"]


def escape_string(input):
    if isinstance(input, str):
        return (
            input.replace(r"{", r"")
            .replace(r"}", r"\}")
            .replace("<", "&lt;")
            .replace(">", "&gt;")
        )
    return input


class Node:
    def __init__(self, tag, children=None):
        self.tag = tag
        if children is None:
            self.children = []
        elif isinstance(children, list):
            self.children = children
        else:
            self.children = [children]

    def add_child(self, child):
        self.children.append(child)

    def __str__(self):
        if len(self.children) == 0:
            return f"<{self.tag}></{self.tag}>"

        if len(self.children) == 1:
            return f"<{self.tag}>{escape_string(self.children[0])}</{self.tag}>"

        output = f"<{self.tag}>\n"
        for child in self.children:
            output += f"{escape_string(child)}\n"
        output += f"</{self.tag}>"
        return output


def wrap_in_code_tag(input):
    return Node("code", input) if input != "" else None


def main():
    # Load the functions from the JSON file
    with open(INPUT_FILENAME, "r") as f:
        raw_data = f.read()
        all_fns = json.loads(raw_data)
    print(f"Loaded {len(all_fns)} functions")

    # Group the functions by `category` field
    category_to_fns = {}
    for fn_name, fn in all_fns.items():
        if fn["category"] not in category_to_fns:
            category_to_fns[fn["category"]] = []
        new_fn = deepcopy(fn)
        new_fn["name"] = fn_name
        category_to_fns[fn["category"]].append(new_fn)
    print(f"Grouped functions into {len(category_to_fns.items())} categories")

    # Write to the output file
    with open(OUTPUT_FILENAME, "w") as f:
        for category, fns in category_to_fns.items():
            f.write(f"# Gofakeit Functions\n\n")
            f.write(f"## {category.capitalize()}\n\n")

            table = Node("table")

            table_header = Node("tr")
            for header in HEADERS:
                table_header.add_child(Node("td", header))
            table.add_child(table_header)

            for fn in fns:
                table_row = Node("tr")
                table_row.add_child(Node("td", wrap_in_code_tag(fn["name"])))
                table_row.add_child(Node("td", fn["description"]))
                table_row.add_child(Node("td", wrap_in_code_tag(fn["output"])))

                table_row_params_col = Node("td")
                if fn["params"] is not None:
                    params_table = Node("table")

                    params_table_header = Node("tr")
                    for header in PARAM_HEADERS:
                        params_table_header.add_child(Node("th", header))
                    params_table.add_child(params_table_header)

                    for param in fn["params"]:
                        params_table_row = Node("tr")
                        params_table_row.add_child(
                            Node("td", wrap_in_code_tag(param["field"]))
                        )
                        params_table_row.add_child(Node("td", param["description"]))
                        params_table_row.add_child(
                            Node("td", wrap_in_code_tag(param["type"]))
                        )
                        params_table_row.add_child(Node("td", param["optional"]))
                        params_table_row.add_child(
                            Node("td", wrap_in_code_tag(param["default"]))
                        )
                        params_table_options_col = Node("td")
                        if param["options"] is not None:
                            for option in param["options"]:
                                params_table_options_col.add_child(
                                    Node("li", wrap_in_code_tag(option))
                                )
                        params_table_row.add_child(params_table_options_col)
                        params_table.add_child(params_table_row)

                    table_row_params_col.add_child(params_table)
                table_row.add_child(table_row_params_col)

                table.add_child(table_row)

            f.write(f"{table}\n\n")


if __name__ == "__main__":
    main()
