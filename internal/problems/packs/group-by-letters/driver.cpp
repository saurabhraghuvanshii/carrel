#include <algorithm>
#include <climits>
#include <cmath>
#include <cstdlib>
#include <exception>
#include <iostream>
#include <list>
#include <map>
#include <queue>
#include <set>
#include <sstream>
#include <stack>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>
using namespace std;

#include "solution.cpp"

// Sorts the words in each group, then the groups, so any order is accepted.
string showGroups(vector<vector<string>> groups) {
    for (auto& g : groups) sort(g.begin(), g.end());
    sort(groups.begin(), groups.end());
    ostringstream o;
    o << '[';
    for (size_t i = 0; i < groups.size(); i++) {
        if (i > 0) o << ", ";
        o << '[';
        for (size_t j = 0; j < groups[i].size(); j++) {
            if (j > 0) o << ", ";
            o << groups[i][j];
        }
        o << ']';
    }
    o << ']';
    return o.str();
}

// Reads T cases from stdin and prints one line per case. Each case is "n" then n words.
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        int n;
        cin >> n;
        vector<string> words(n);
        for (auto& w : words) cin >> w;
        try {
            out << showGroups(groupByLetters(words)) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
