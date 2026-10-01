#include <algorithm>
#include <climits>
#include <cmath>
#include <cstdlib>
#include <exception>
#include <iomanip>
#include <iostream>
#include <list>
#include <map>
#include <queue>
#include <set>
#include <sstream>
#include <stdexcept>
#include <stack>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>
using namespace std;

struct Node {
    int val;
    vector<Node*> neighbors;
    Node(int v) : val(v) {}
};

#include "solution.cpp"

// Reads m edges of width numbers each: "u v" or "u v w".
vector<vector<int>> readEdges(int m, int width) {
    vector<vector<int>> edges(m, vector<int>(width));
    for (auto& e : edges)
        for (auto& x : e) cin >> x;
    return edges;
}

// Reads T cases from stdin and prints one line per case. Each case is "n m", then m edges "u v".
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        int n, m;
        cin >> n >> m;
        vector<vector<int>> edges = readEdges(m, 2);
        vector<Node*> nodes;
        for (int i = 0; i < n; i++) nodes.push_back(new Node(i));
        for (auto& e : edges) {
            nodes[e[0]]->neighbors.push_back(nodes[e[1]]);
            nodes[e[1]]->neighbors.push_back(nodes[e[0]]);
        }
        try {
            Node* copy = cloneGraph(nodes[0]);
            unordered_set<Node*> originals(nodes.begin(), nodes.end()), seen;
            vector<Node*> byVal(n, nullptr);
            queue<Node*> q;
            string problem = copy == nullptr ? "nothing was returned" : "";
            if (copy) {
                q.push(copy);
                seen.insert(copy);
            }
            while (problem.empty() && !q.empty()) {
                Node* x = q.front();
                q.pop();
                if (originals.count(x)) {
                    problem = "the copy uses original node " + to_string(x->val);
                } else if (x->val < 0 || x->val >= n || byVal[x->val]) {
                    problem = "node " + to_string(x->val) + " is copied more than once or has a wrong value";
                } else {
                    byVal[x->val] = x;
                    for (Node* y : x->neighbors)
                        if (y && seen.insert(y).second) q.push(y);
                }
            }
            for (int i = 0; problem.empty() && i < n; i++)
                if (!byVal[i]) problem = "node " + to_string(i) + " is missing from the copy";
            if (!problem.empty()) {
                out << "not a copy: " << problem << '\n';
            } else {
                out << '[';
                for (int i = 0; i < n; i++) {
                    out << (i ? ", [" : "[");
                    for (size_t j = 0; j < byVal[i]->neighbors.size(); j++) out << (j ? ", " : "") << byVal[i]->neighbors[j]->val;
                    out << ']';
                }
                out << "]\n";
            }
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
