#include <algorithm>
#include <climits>
#include <cmath>
#include <cstdlib>
#include <exception>
#include <iostream>
#include <map>
#include <queue>
#include <set>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>
using namespace std;

#include "solution.cpp"

// Reads T cases from stdin, calls the learner's code, prints one line per case.
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
        vector<vector<int>> drivers(n, vector<int>(4));
        for (auto& d : drivers) {
            for (auto& v : d) cin >> v;
        }
        int rx, ry;
        cin >> rx >> ry;
        try {
            out << nearestDriver(drivers, rx, ry) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
    }
    out.flush();
    return 0;
}
