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

// Reads T cases from stdin and prints one line per case. Each case is the zone "x1 y1 x2 y2", then "n" and n lines "x y".
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        vector<int> zone(4);
        for (auto& x : zone) cin >> x;
        int n;
        cin >> n;
        vector<vector<int>> path(n, vector<int>(2));
        for (auto& r : path)
            for (auto& x : r) cin >> x;
        try {
            out << entries(zone, path) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
