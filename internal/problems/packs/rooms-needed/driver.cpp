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

vector<int> readArray() {
    int n;
    cin >> n;
    vector<int> a(n);
    for (auto& x : a) cin >> x;
    return a;
}

string show(const vector<int>& a) {
    ostringstream o;
    o << '[';
    for (size_t i = 0; i < a.size(); i++) {
        if (i > 0) o << ", ";
        o << a[i];
    }
    o << ']';
    return o.str();
}

vector<vector<int>> readGrid() {
    int rows, cols;
    cin >> rows >> cols;
    vector<vector<int>> g(rows, vector<int>(cols));
    for (auto& row : g)
        for (auto& x : row) cin >> x;
    return g;
}

string show(const vector<vector<int>>& g) {
    ostringstream o;
    o << '[';
    for (size_t r = 0; r < g.size(); r++) {
        if (r > 0) o << ", ";
        o << show(g[r]);
    }
    o << ']';
    return o.str();
}

// Reads "n" then n pairs "start end".
vector<vector<int>> readIntervals() {
    int n;
    cin >> n;
    vector<vector<int>> iv(n, vector<int>(2));
    for (auto& p : iv) cin >> p[0] >> p[1];
    return iv;
}

// Reads T cases from stdin and prints one line per case. Each case is "n" then n lines "start end".
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        vector<vector<int>> meetings = readIntervals();
        try {
            out << roomsNeeded(meetings) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
