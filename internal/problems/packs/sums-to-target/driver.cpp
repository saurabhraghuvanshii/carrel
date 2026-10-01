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

// Sorts each triple, then the list, so any order is accepted.
string showTriples(vector<vector<int>> lists) {
    for (auto& l : lists) sort(l.begin(), l.end());
    sort(lists.begin(), lists.end());
    ostringstream o;
    o << '[';
    for (size_t i = 0; i < lists.size(); i++) {
        if (i > 0) o << ", ";
        o << '[';
        for (size_t j = 0; j < lists[i].size(); j++) {
            if (j > 0) o << ", ";
            o << lists[i][j];
        }
        o << ']';
    }
    o << ']';
    return o.str();
}

// Reads T cases from stdin and prints one line per case. Each case is "n", n candidates, then the target.
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        vector<int> candidates = readArray();
        int target;
        cin >> target;
        try {
            out << showTriples(combinations(candidates, target)) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
