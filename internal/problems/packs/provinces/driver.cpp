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
        try {
            out << countGroups(n, edges) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
