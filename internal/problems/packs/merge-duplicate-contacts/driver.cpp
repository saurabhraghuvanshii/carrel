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

// Reads T cases from stdin and prints one line per case. Each case is "n", then n lines "k email ... email".
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
        vector<vector<string>> contacts(n);
        for (auto& one : contacts) {
            int count;
            cin >> count;
            one.resize(count);
            for (auto& address : one) cin >> address;
        }
        try {
            out << show(keepers(contacts)) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
