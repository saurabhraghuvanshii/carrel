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

// Reads T cases from stdin and prints one line per case. Each case is "n", then n lines "start end", then the new meeting "start end".
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
        vector<vector<int>> meetings(n, vector<int>(2));
        for (auto& r : meetings)
            for (auto& x : r) cin >> x;
        int start, end;
        cin >> start >> end;
        try {
            out << clashes(meetings, start, end) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
