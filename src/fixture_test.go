// Tests for the behavior of the test fixture system.

package gocheck_test


import (
    "gocheck"
    "regexp"
)


// -----------------------------------------------------------------------
// Fixture test suite.

type FixtureS struct{}

var fixtureS = gocheck.Suite(&FixtureS{})

func (s *FixtureS) TestCountSuite(t *gocheck.T) {
    suitesRun += 1
}


// -----------------------------------------------------------------------
// Basic fixture ordering verification.

func (s *FixtureS) TestOrder(t *gocheck.T) {
    helper := FixtureHelper{}
    gocheck.Run(&helper, nil)
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "Test1")
    t.CheckEqual(helper.calls[3], "TearDownTest")
    t.CheckEqual(helper.calls[4], "SetUpTest")
    t.CheckEqual(helper.calls[5], "Test2")
    t.CheckEqual(helper.calls[6], "TearDownTest")
    t.CheckEqual(helper.calls[7], "TearDownSuite")
    t.CheckEqual(helper.n, 8)
}


// -----------------------------------------------------------------------
// Check the behavior when panics occur within tests and fixtures.

func (s *FixtureS) TestPanicOnTest(t *gocheck.T) {
    helper := FixtureHelper{panicOn: "Test1"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "Test1")
    t.CheckEqual(helper.calls[3], "TearDownTest")
    t.CheckEqual(helper.calls[4], "SetUpTest")
    t.CheckEqual(helper.calls[5], "Test2")
    t.CheckEqual(helper.calls[6], "TearDownTest")
    t.CheckEqual(helper.calls[7], "TearDownSuite")
    t.CheckEqual(helper.n, 8)

    expected := "^\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: FixtureHelper.Test1\n\n" +
                "\\.\\.\\. Panic: Test1 \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.Test1\n$"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        t.Error("Bad expression:", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnSetUpTest(t *gocheck.T) {
    helper := FixtureHelper{panicOn: "SetUpTest"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "TearDownTest")
    t.CheckEqual(helper.calls[3], "TearDownSuite")
    t.CheckEqual(helper.n, 4)

    expected := "^\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper\\.SetUpTest\n\n" +
                "\\.\\.\\. Panic: SetUpTest \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.SetUpTest\n" +
                "\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: Fixture has panicked " +
                "\\(see related PANIC\\)\n$"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        t.Error("Bad expression:", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnTearDownTest(t *gocheck.T) {
    helper := FixtureHelper{panicOn: "TearDownTest"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "Test1")
    t.CheckEqual(helper.calls[3], "TearDownTest")
    t.CheckEqual(helper.calls[4], "TearDownSuite")
    t.CheckEqual(helper.n, 5)

    expected := "^\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper.TearDownTest\n\n" +
                "\\.\\.\\. Panic: TearDownTest \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.TearDownTest\n" +
                "\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: Fixture has panicked " +
                "\\(see related PANIC\\)\n$"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        t.Error("Bad expression:", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnSetUpSuite(t *gocheck.T) {
    helper := FixtureHelper{panicOn: "SetUpSuite"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "TearDownSuite")
    t.CheckEqual(helper.n, 2)

    expected := "^\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper.SetUpSuite\n\n" +
                "\\.\\.\\. Panic: SetUpSuite \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.SetUpSuite\n$"

    // XXX Changing the expression above to not match breaks Go. WTF?

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        t.Error("Bad expression:", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnTearDownSuite(t *gocheck.T) {
    helper := FixtureHelper{panicOn: "TearDownSuite"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "Test1")
    t.CheckEqual(helper.calls[3], "TearDownTest")
    t.CheckEqual(helper.calls[4], "SetUpTest")
    t.CheckEqual(helper.calls[5], "Test2")
    t.CheckEqual(helper.calls[6], "TearDownTest")
    t.CheckEqual(helper.calls[7], "TearDownSuite")
    t.CheckEqual(helper.n, 8)

    expected := "^\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper.TearDownSuite\n\n" +
                "\\.\\.\\. Panic: TearDownSuite \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.TearDownSuite\n$"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        t.Error("Bad expression:", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}


// -----------------------------------------------------------------------
// A wrong argument on a test or fixture will produce a nice error.

func (s *FixtureS) TestPanicOnWrongTestArg(t *gocheck.T) {
    helper := WrongTestArgHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "TearDownTest")
    t.CheckEqual(helper.calls[3], "SetUpTest")
    t.CheckEqual(helper.calls[4], "Test2")
    t.CheckEqual(helper.calls[5], "TearDownTest")
    t.CheckEqual(helper.calls[6], "TearDownSuite")
    t.CheckEqual(helper.n, 7)

    expected := "^\n-+\n" +
                "PANIC: fixture_test\\.go:[0-9]+: " +
                "WrongTestArgHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: WrongTestArgHelper\\.Test1 argument " +
                "should be \\*gocheck\\.T\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpTestArg(t *gocheck.T) {
    helper := WrongSetUpTestArgHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpTestArgHelper\\.SetUpTest\n\n" +
        "\\.\\.\\. Panic: WrongSetUpTestArgHelper\\.SetUpTest argument " +
        "should be \\*gocheck\\.F\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpSuiteArg(t *gocheck.T) {
    helper := WrongSetUpSuiteArgHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpSuiteArgHelper\\.SetUpSuite\n\n" +
        "\\.\\.\\. Panic: WrongSetUpSuiteArgHelper\\.SetUpSuite argument " +
        "should be \\*gocheck\\.F\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}


// -----------------------------------------------------------------------
// Nice errors also when tests or fixture have wrong arg count.

func (s *FixtureS) TestPanicOnWrongTestArgCount(t *gocheck.T) {
    helper := WrongTestArgCountHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "TearDownTest")
    t.CheckEqual(helper.calls[3], "SetUpTest")
    t.CheckEqual(helper.calls[4], "Test2")
    t.CheckEqual(helper.calls[5], "TearDownTest")
    t.CheckEqual(helper.calls[6], "TearDownSuite")
    t.CheckEqual(helper.n, 7)

    expected := "^\n-+\n" +
                "PANIC: fixture_test\\.go:[0-9]+: " +
                "WrongTestArgCountHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: WrongTestArgCountHelper\\.Test1 argument " +
                "should be \\*gocheck\\.T\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpTestArgCount(t *gocheck.T) {
    helper := WrongSetUpTestArgCountHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpTestArgCountHelper\\.SetUpTest\n\n" +
        "\\.\\.\\. Panic: WrongSetUpTestArgCountHelper\\.SetUpTest argument " +
        "should be \\*gocheck\\.F\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpSuiteArgCount(t *gocheck.T) {
    helper := WrongSetUpSuiteArgCountHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpSuiteArgCountHelper\\.SetUpSuite\n\n" +
        "\\.\\.\\. Panic: WrongSetUpSuiteArgCountHelper" +
        "\\.SetUpSuite argument should be \\*gocheck\\.F\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}


// -----------------------------------------------------------------------
// Helper test suites with wrong function arguments.

type WrongTestArgHelper struct {
    FixtureHelper
}

func (s *WrongTestArgHelper) Test1(t int) {
}

// ----

type WrongSetUpTestArgHelper struct {
    FixtureHelper
}

func (s *WrongSetUpTestArgHelper) SetUpTest(t int) {
}

type WrongSetUpSuiteArgHelper struct {
    FixtureHelper
}

func (s *WrongSetUpSuiteArgHelper) SetUpSuite(t int) {
}

type WrongTestArgCountHelper struct {
    FixtureHelper
}

func (s *WrongTestArgCountHelper) Test1(t *gocheck.T, i int) {
}

type WrongSetUpTestArgCountHelper struct {
    FixtureHelper
}

func (s *WrongSetUpTestArgCountHelper) SetUpTest(f *gocheck.F, i int) {
}

type WrongSetUpSuiteArgCountHelper struct {
    FixtureHelper
}

func (s *WrongSetUpSuiteArgCountHelper) SetUpSuite(f *gocheck.F, i int) {
}


// -----------------------------------------------------------------------
// Ensure fixture doesn't without tests.

type NoTestsHelper struct{
    hasRun bool
}

func (s *NoTestsHelper) SetUpSuite(f *gocheck.F) {
    s.hasRun = true
}

func (s *NoTestsHelper) TearDownSuite(f *gocheck.F) {
    s.hasRun = true
}

func (s *FixtureS) TestFixtureDoesntRunWithoutTests(t *gocheck.T) {
    helper := NoTestsHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(helper.hasRun, false)
}
